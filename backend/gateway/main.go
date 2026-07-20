package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var jwtKey = []byte("super_secret_key_123") // Needs to match auth-service
var mqttClient mqtt.Client

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initMQTT() {
	opts := mqtt.NewClientOptions()
	mqttBroker := getEnvOrDefault("MQTT_BROKER", "localhost:1883")
	opts.AddBroker(fmt.Sprintf("tcp://%s", mqttBroker))
	opts.SetClientID("api_gateway")

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("MQTT Error in Gateway: %v", token.Error())
	} else {
		log.Println("Gateway Connected to MQTT Broker for Commands")
	}
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// fallback to query param for websockets
			tokenString = c.Query("token")
		} else {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for MVP
	},
}

func proxyRequest(method string, url string, c *gin.Context) {
	var bodyReader io.Reader
	if c.Request.Body != nil {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}
	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

func wsHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer ws.Close()

	for {
		time.Sleep(3 * time.Second)
		
		resp, err := http.Get("http://localhost:8082/latest")
		if err != nil { continue }
		body, _ := io.ReadAll(resp.Body)
		var latest map[string]float64
		json.Unmarshal(body, &latest)
		resp.Body.Close()
		
		resp2, err := http.Get("http://localhost:8082/alerts")
		var alerts []interface{}
		if err == nil {
			body2, _ := io.ReadAll(resp2.Body)
			json.Unmarshal(body2, &alerts)
			resp2.Body.Close()
		}

		msg := map[string]interface{}{
			"type": "update",
			"data": map[string]interface{}{
				"latest": latest,
				"alerts": alerts,
			},
		}

		if err := ws.WriteJSON(msg); err != nil {
			break
		}
	}
}

func main() {
	initMQTT()
	r := gin.Default()
	
	// Basic CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public Auth endpoints
	authSvcPub := getEnvOrDefault("AUTH_SVC_URL", "http://localhost:8083")
	r.POST("/api/auth/login", func(c *gin.Context) {
		proxyRequest("POST", authSvcPub+"/login", c)
	})
	r.POST("/api/auth/register", func(c *gin.Context) {
		proxyRequest("POST", authSvcPub+"/register", c)
	})

	// Protected endpoints
	protected := r.Group("/")
	protected.Use(JWTMiddleware())
	
	authSvc := getEnvOrDefault("AUTH_SVC_URL", "http://localhost:8083")
	deviceSvc := getEnvOrDefault("DEVICE_SVC_URL", "http://localhost:8081")
	telemetrySvc := getEnvOrDefault("TELEMETRY_SVC_URL", "http://localhost:8082")
	ruleSvc := getEnvOrDefault("RULE_SVC_URL", "http://localhost:8082")

	protected.GET("/api/devices", func(c *gin.Context) { proxyRequest("GET", deviceSvc+"/devices", c) })
	protected.POST("/api/devices", func(c *gin.Context) { proxyRequest("POST", deviceSvc+"/devices", c) })
	protected.DELETE("/api/devices/:id", func(c *gin.Context) { proxyRequest("DELETE", deviceSvc+"/devices/"+c.Param("id"), c) })
	
	protected.GET("/api/history", func(c *gin.Context) {
		query := c.Request.URL.RawQuery
		url := telemetrySvc + "/history"
		if query != "" {
			url += "?" + query
		}
		proxyRequest("GET", url, c)
	})
	protected.GET("/api/alerts", func(c *gin.Context) { proxyRequest("GET", telemetrySvc+"/alerts", c) })
	
	protected.GET("/api/rules", func(c *gin.Context) { proxyRequest("GET", ruleSvc+"/rules", c) })
	protected.POST("/api/rules", func(c *gin.Context) { proxyRequest("POST", ruleSvc+"/rules", c) })
	protected.DELETE("/api/rules/:id", func(c *gin.Context) { proxyRequest("DELETE", ruleSvc+"/rules/"+c.Param("id"), c) })

	protected.POST("/api/devices/:id/command", func(c *gin.Context) {
		id := c.Param("id")
		var payload map[string]string
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		
		cmd, exists := payload["command"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "command required"})
			return
		}

		topic := fmt.Sprintf("iot/sensors/%s/command", id)
		if mqttClient != nil && mqttClient.IsConnected() {
			mqttClient.Publish(topic, 1, false, cmd)
			c.JSON(http.StatusOK, gin.H{"status": "command sent"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MQTT disconnected"})
		}
	})

	protected.GET("/ws", wsHandler)

	log.Println("API Gateway running on port 8080")
	r.Run(":8080")
}
