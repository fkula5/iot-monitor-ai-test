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

func createClient(deviceID string) mqtt.Client {
	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		mqttBroker = "localhost:1883"
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", mqttBroker))
	opts.SetClientID(deviceID)
	
	// Set Last Will and Testament (LWT)
	statusTopic := fmt.Sprintf("iot/sensors/%s/status", deviceID)
	opts.SetWill(statusTopic, "offline", 1, true)
	
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker for %s: %v", deviceID, token.Error())
	}
	
	// Publish Online status immediately after connection
	client.Publish(statusTopic, 1, true, "online")
	return client
}

func simulateDevice(deviceID string, baseValue float64, variance float64) {
	client := createClient(deviceID)
	defer client.Disconnect(250)

	active := true

	// Subscribe to commands
	commandTopic := fmt.Sprintf("iot/sensors/%s/command", deviceID)
	client.Subscribe(commandTopic, 1, func(c mqtt.Client, m mqtt.Message) {
		cmd := string(m.Payload())
		log.Printf("[%s] Received command: %s", deviceID, cmd)
		if cmd == "TURN_OFF" {
			active = false
			client.Publish(fmt.Sprintf("iot/sensors/%s/status", deviceID), 1, true, "offline")
		} else if cmd == "TURN_ON" {
			active = true
			client.Publish(fmt.Sprintf("iot/sensors/%s/status", deviceID), 1, true, "online")
		} else if cmd == "RESTART" {
			log.Printf("[%s] Restarting...", deviceID)
			active = false
			client.Publish(fmt.Sprintf("iot/sensors/%s/status", deviceID), 1, true, "offline")
			time.Sleep(2 * time.Second)
			active = true
			client.Publish(fmt.Sprintf("iot/sensors/%s/status", deviceID), 1, true, "online")
			log.Printf("[%s] Restart complete.", deviceID)
		}
	})

	currentValue := baseValue
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !active {
			continue
		}
		
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
	log.Println("Starting Edge Simulator...")

	go simulateDevice("dev-1", 22.0, 1.0)
	go simulateDevice("dev-2", 45.0, 2.0)
	go simulateDevice("dev-4", -4.0, 1.0)
	
	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	
	log.Println("Shutting down Edge Simulator...")
}
