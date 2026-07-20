package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TelemetryPayload struct {
	DeviceID string  `json:"deviceId"`
	Value    float64 `json:"value"`
}

func createClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("edge_simulator")
	
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}
	return client
}

func simulateDevice(client mqtt.Client, deviceID string, baseValue float64, variance float64) {
	currentValue := baseValue
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentValue = currentValue + (rand.Float64()-0.5)*variance
		if currentValue < -20 { currentValue = -20 } // Floor
		
		payload := TelemetryPayload{
			DeviceID: deviceID,
			Value:    currentValue,
		}
		
		jsonData, _ := json.Marshal(payload)
		topic := fmt.Sprintf("iot/sensors/%s/telemetry", deviceID)
		
		token := client.Publish(topic, 0, false, jsonData)
		token.Wait()
		
		log.Printf("Published to %s: %s", topic, string(jsonData))
	}
}

func main() {
	client := createClient()
	defer client.Disconnect(250)
	
	log.Println("Edge Simulator Connected to MQTT Broker. Starting simulation...")

	go simulateDevice(client, "dev-1", 22.0, 1.0)
	go simulateDevice(client, "dev-2", 45.0, 2.0)
	go simulateDevice(client, "dev-4", -4.0, 1.0)
	
	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	
	log.Println("Shutting down Edge Simulator...")
}
