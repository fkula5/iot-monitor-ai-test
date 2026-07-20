package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type DataPoint struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

type Alert struct {
	ID        int64  `json:"id"`
	DeviceID  string `json:"deviceId"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

type TelemetryPayload struct {
	DeviceID string  `json:"deviceId"`
	Value    float64 `json:"value"`
}

var (
	mu           sync.RWMutex
	alerts       []Alert
	latestValues = map[string]float64{}
	influxClient influxdb2.Client
	influxURL    = "http://localhost:8086"
	influxToken  = "my-super-secret-auth-token"
	influxOrg    = "iot_org"
	influxBucket = "telemetry_bucket"
)

func initInflux() {
	influxClient = influxdb2.NewClient(influxURL, influxToken)
}

func initMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("telemetry_service")

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		var payload TelemetryPayload
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Printf("Invalid MQTT payload: %v", err)
			return
		}

		mu.Lock()
		latestValues[payload.DeviceID] = payload.Value
		
		if payload.Value > 40.0 {
			alert := Alert{
				ID:        time.Now().UnixNano(),
				DeviceID:  payload.DeviceID,
				Message:   "Krytyczna temperatura przekroczyła 40°C! (MQTT/Alert Engine)",
				Type:      "danger",
				Timestamp: time.Now().Format(time.RFC3339),
			}
			alerts = append([]Alert{alert}, alerts...)
			if len(alerts) > 10 {
				alerts = alerts[:10]
			}
		}
		mu.Unlock()

		writeAPI := influxClient.WriteAPIBlocking(influxOrg, influxBucket)
		p := influxdb2.NewPoint("temperature",
			map[string]string{"device": payload.DeviceID},
			map[string]interface{}{"value": payload.Value},
			time.Now())
		writeAPI.WritePoint(context.Background(), p)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT: %v", token.Error())
	}

	if token := client.Subscribe("iot/sensors/+/telemetry", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to MQTT: %v", token.Error())
	}
	
	log.Println("Subscribed to MQTT topics.")
}

func main() {
	initInflux()
	initMQTT()

	r := gin.Default()

	r.GET("/history", func(c *gin.Context) {
		timeRange := c.DefaultQuery("range", "-15m")
		
		if !strings.HasPrefix(timeRange, "-") {
			timeRange = "-15m"
		}

		queryAPI := influxClient.QueryAPI(influxOrg)
		
		var window = "5s"
		if timeRange == "-1h" {
			window = "30s"
		} else if timeRange == "-24h" {
			window = "5m"
		}

		query := `from(bucket: "` + influxBucket + `") 
			|> range(start: ` + timeRange + `) 
			|> filter(fn: (r) => r._measurement == "temperature") 
			|> aggregateWindow(every: ` + window + `, fn: mean, createEmpty: false)
			|> keep(columns: ["_time", "_value", "device"])`
		
		result, err := queryAPI.Query(context.Background(), query)
		
		history := map[string][]DataPoint{}
		if err == nil {
			for result.Next() {
				dev := result.Record().ValueByKey("device").(string)
				val := result.Record().Value().(float64)
				t := result.Record().Time().Format("15:04:05")
				history[dev] = append(history[dev], DataPoint{Time: t, Value: val})
			}
		}
		
		c.JSON(http.StatusOK, history)
	})
	
	r.GET("/latest", func(c *gin.Context) {
		mu.RLock()
		defer mu.RUnlock()
		c.JSON(http.StatusOK, latestValues)
	})

	r.GET("/alerts", func(c *gin.Context) {
		mu.RLock()
		defer mu.RUnlock()
		c.JSON(http.StatusOK, alerts)
	})

	log.Println("Telemetry Service running on port 8082")
	r.Run(":8082")
}
