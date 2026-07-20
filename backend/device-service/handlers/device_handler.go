package handlers

import (
	"device-service/models"
	"device-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	repo repository.DeviceRepository
}

func NewDeviceHandler(repo repository.DeviceRepository) *DeviceHandler {
	return &DeviceHandler{repo: repo}
}

func (h *DeviceHandler) GetDevices(c *gin.Context) {
	devices, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var newDevice models.Device
	if err := c.BindJSON(&newDevice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	newDevice.Status = "online"
	newDevice.Battery = 100
	newDevice.Uptime = "0 dni 0h"
	if err := h.repo.Create(&newDevice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
		return
	}
	c.JSON(http.StatusCreated, newDevice)
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
