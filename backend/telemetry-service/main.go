package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

var (
	mu          sync.RWMutex
	alerts      []Alert
	latestValues = map[string]float64{
		"dev-1": 22.5,
		"dev-2": 45.0,
		"dev-4": -4.2,
	}
	influxClient influxdb2.Client
	influxURL    = "http://localhost:8086"
	influxToken  = "my-super-secret-auth-token"
	influxOrg    = "iot_org"
	influxBucket = "telemetry_bucket"
)

func initInflux() {
	influxClient = influxdb2.NewClient(influxURL, influxToken)
}

func generateHistoryToInflux() {
	writeAPI := influxClient.WriteAPIBlocking(influxOrg, influxBucket)
	now := time.Now()
	for i := 15; i >= 0; i-- {
		t := now.Add(time.Duration(-i*3) * time.Second)
		p1 := influxdb2.NewPoint("temperature",
			map[string]string{"device": "dev-1"},
			map[string]interface{}{"value": 22.5 + (rand.Float64()-0.5)*2},
			t)
		p4 := influxdb2.NewPoint("temperature",
			map[string]string{"device": "dev-4"},
			map[string]interface{}{"value": -4.2 + (rand.Float64()-0.5)*1},
			t)
		writeAPI.WritePoint(context.Background(), p1)
		writeAPI.WritePoint(context.Background(), p4)
	}
}

func simulator() {
	writeAPI := influxClient.WriteAPIBlocking(influxOrg, influxBucket)
	for {
		time.Sleep(3 * time.Second)
		mu.Lock()
		
		val1 := latestValues["dev-1"] + (rand.Float64()-0.5)*1.0
		val4 := latestValues["dev-4"] + (rand.Float64()-0.5)*1.0
		if val4 > 0 { val4 = -1.0 }

		latestValues["dev-1"] = val1
		latestValues["dev-4"] = val4

		// Write to InfluxDB
		t := time.Now()
		p1 := influxdb2.NewPoint("temperature", map[string]string{"device": "dev-1"}, map[string]interface{}{"value": val1}, t)
		p4 := influxdb2.NewPoint("temperature", map[string]string{"device": "dev-4"}, map[string]interface{}{"value": val4}, t)
		writeAPI.WritePoint(context.Background(), p1)
		writeAPI.WritePoint(context.Background(), p4)

		if rand.Float64() > 0.9 {
			alert := Alert{
				ID:        time.Now().UnixNano(),
				DeviceID:  "dev-1",
				Message:   "Odnotowano nagły skok parametru (InfluxDB Backend).",
				Type:      "warning",
				Timestamp: time.Now().Format(time.RFC3339),
			}
			alerts = append([]Alert{alert}, alerts...)
			if len(alerts) > 10 {
				alerts = alerts[:10]
			}
		}
		mu.Unlock()
	}
}

func main() {
	initInflux()
	generateHistoryToInflux()
	go simulator()

	r := gin.Default()

	r.GET("/history", func(c *gin.Context) {
		queryAPI := influxClient.QueryAPI(influxOrg)
		query := `from(bucket: "` + influxBucket + `") 
			|> range(start: -10m) 
			|> filter(fn: (r) => r._measurement == "temperature") 
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
