package main

import (
	"log"
	
	"device-service/handlers"
	"device-service/models"
	"device-service/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupRouter(repo repository.DeviceRepository) *gin.Engine {
	r := gin.Default()
	handler := handlers.NewDeviceHandler(repo)
	
	r.GET("/devices", handler.GetDevices)
	r.POST("/devices", handler.CreateDevice)
	r.DELETE("/devices/:id", handler.DeleteDevice)
	
	return r
}

func main() {
	dsn := "host=localhost user=admin password=adminpassword dbname=iot_db port=5432 sslmode=disable TimeZone=Europe/Warsaw"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&models.Device{})

	repo := repository.NewDeviceRepository(db)

	count, _ := repo.Count()
	if count == 0 {
		repo.Create(&models.Device{ID: "dev-1", Name: "Czujnik Temp. - Magazyn A", Type: "temperature", Status: "online", Battery: 85, Uptime: "14 dni 2h", Unit: "°C"})
		repo.Create(&models.Device{ID: "dev-2", Name: "Czujnik Wilgotności - Magazyn A", Type: "humidity", Status: "online", Battery: 90, Uptime: "14 dni 2h", Unit: "%"})
		repo.Create(&models.Device{ID: "dev-3", Name: "Stacja Pogodowa - Zewnętrzna", Type: "weather", Status: "offline", Battery: 0, Uptime: "-", Unit: ""})
		repo.Create(&models.Device{ID: "dev-4", Name: "Chłodnia - Sektor B", Type: "temperature", Status: "online", Battery: 100, Uptime: "5 dni 10h", Unit: "°C"})
	}

	r := SetupRouter(repo)

	log.Println("Device Service running on port 8081")
	r.Run(":8081")
}
