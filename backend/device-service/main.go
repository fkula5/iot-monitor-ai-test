package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Device struct {
	ID      string  `json:"id" gorm:"primaryKey"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Status  string  `json:"status"`
	Battery int     `json:"battery"`
	Uptime  string  `json:"uptime"`
	Unit    string  `json:"unit"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("devices.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&Device{})

	// Seed data
	var count int64
	db.Model(&Device{}).Count(&count)
	if count == 0 {
		db.Create([]Device{
			{ID: "dev-1", Name: "Czujnik Temp. - Magazyn A", Type: "temperature", Status: "online", Battery: 85, Uptime: "14 dni 2h", Unit: "°C"},
			{ID: "dev-2", Name: "Czujnik Wilgotności - Magazyn A", Type: "humidity", Status: "online", Battery: 90, Uptime: "14 dni 2h", Unit: "%"},
			{ID: "dev-3", Name: "Stacja Pogodowa - Zewnętrzna", Type: "weather", Status: "offline", Battery: 0, Uptime: "-", Unit: ""},
			{ID: "dev-4", Name: "Chłodnia - Sektor B", Type: "temperature", Status: "online", Battery: 100, Uptime: "5 dni 10h", Unit: "°C"},
		})
	}

	r := gin.Default()

	r.GET("/devices", func(c *gin.Context) {
		var devices []Device
		db.Find(&devices)
		c.JSON(http.StatusOK, devices)
	})
	
	r.POST("/devices", func(c *gin.Context) {
		var newDevice Device
		if err := c.BindJSON(&newDevice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		newDevice.Status = "online"
		newDevice.Battery = 100
		newDevice.Uptime = "0 dni 0h"
		if err := db.Create(&newDevice).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
			return
		}
		c.JSON(http.StatusCreated, newDevice)
	})

	r.DELETE("/devices/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Device{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
	})

	log.Println("Device Service running on port 8081")
	r.Run(":8081")
}
