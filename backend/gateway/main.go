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
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
			if canWriteDevices, ok := claims["canWriteDevices"].(bool); ok {
				c.Set("canWriteDevices", canWriteDevices)
			}
			if canWriteRules, ok := claims["canWriteRules"].(bool); ok {
				c.Set("canWriteRules", canWriteRules)
			}
			if canManageUsers, ok := claims["canManageUsers"].(bool); ok {
				c.Set("canManageUsers", canManageUsers)
			}
		}

		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(permission)
		if !exists || val != true {
			// Fallback: Admin always has access to everything
			role, roleExists := c.Get("role")
			if roleExists && role == "Admin" {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Permission %s required", permission)})
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

	telemetrySvc := getEnvOrDefault("TELEMETRY_SVC_URL", "http://localhost:8082")
	deviceSvc := getEnvOrDefault("DEVICE_SVC_URL", "http://localhost:8081")
	
	for {
		time.Sleep(3 * time.Second)
		
		resp, err := http.Get(telemetrySvc + "/latest")
		if err != nil { continue }
		body, _ := io.ReadAll(resp.Body)
		var latest map[string]float64
		json.Unmarshal(body, &latest)
		resp.Body.Close()
		
		resp2, err := http.Get(telemetrySvc + "/alerts")
		var alerts []interface{}
		if err == nil {
			body2, _ := io.ReadAll(resp2.Body)
			json.Unmarshal(body2, &alerts)
			resp2.Body.Close()
		}

		resp3, err := http.Get(deviceSvc + "/devices")
		var devices []interface{}
		if err == nil {
			body3, _ := io.ReadAll(resp3.Body)
			json.Unmarshal(body3, &devices)
			resp3.Body.Close()
		}

		msg := map[string]interface{}{
			"type": "update",
			"data": map[string]interface{}{
				"latest": latest,
				"alerts": alerts,
				"devices": devices,
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
	
	deviceSvc := getEnvOrDefault("DEVICE_SVC_URL", "http://localhost:8081")
	telemetrySvc := getEnvOrDefault("TELEMETRY_SVC_URL", "http://localhost:8082")
	ruleSvc := getEnvOrDefault("RULE_SVC_URL", "http://localhost:8082")

	protected.GET("/api/devices", func(c *gin.Context) { proxyRequest("GET", deviceSvc+"/devices", c) })
	protected.POST("/api/devices", RequirePermission("canWriteDevices"), func(c *gin.Context) { proxyRequest("POST", deviceSvc+"/devices", c) })
	protected.DELETE("/api/devices/:id", RequirePermission("canWriteDevices"), func(c *gin.Context) { proxyRequest("DELETE", deviceSvc+"/devices/"+c.Param("id"), c) })
	
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
	protected.POST("/api/rules", RequirePermission("canWriteRules"), func(c *gin.Context) { proxyRequest("POST", ruleSvc+"/rules", c) })
	protected.DELETE("/api/rules/:id", RequirePermission("canWriteRules"), func(c *gin.Context) { proxyRequest("DELETE", ruleSvc+"/rules/"+c.Param("id"), c) })

	authSvcPriv := getEnvOrDefault("AUTH_SVC_URL", "http://localhost:8083")
	protected.GET("/api/users", RequirePermission("canManageUsers"), func(c *gin.Context) { proxyRequest("GET", authSvcPriv+"/users", c) })
	protected.PUT("/api/users/:id/role", RequirePermission("canManageUsers"), func(c *gin.Context) { proxyRequest("PUT", authSvcPriv+"/users/"+c.Param("id")+"/role", c) })
	protected.DELETE("/api/users/:id", RequirePermission("canManageUsers"), func(c *gin.Context) { proxyRequest("DELETE", authSvcPriv+"/users/"+c.Param("id"), c) })

	protected.GET("/api/roles", RequirePermission("canManageUsers"), func(c *gin.Context) { proxyRequest("GET", authSvcPriv+"/roles", c) })
	protected.POST("/api/roles", RequirePermission("canManageUsers"), func(c *gin.Context) { proxyRequest("POST", authSvcPriv+"/roles", c) })
	protected.DELETE("/api/roles/:name", RequirePermission("canManageUsers"), func(c *gin.Context) { proxyRequest("DELETE", authSvcPriv+"/roles/"+c.Param("name"), c) })

	protected.POST("/api/devices/:id/command", RequirePermission("canWriteDevices"), func(c *gin.Context) {
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
