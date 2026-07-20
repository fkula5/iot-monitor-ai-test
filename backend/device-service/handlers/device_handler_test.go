package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"device-service/handlers"
	"device-service/models"
	"device-service/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRepo struct {
	devices []models.Device
}

func (m *MockRepo) GetAll() ([]models.Device, error) {
	return m.devices, nil
}
func (m *MockRepo) Create(d *models.Device) error {
	m.devices = append(m.devices, *d)
	return nil
}
func (m *MockRepo) Delete(id string) error {
	var remaining []models.Device
	for _, d := range m.devices {
		if d.ID != id {
			remaining = append(remaining, d)
		}
	}
	m.devices = remaining
	return nil
}
func (m *MockRepo) Count() (int64, error) {
	return int64(len(m.devices)), nil
}

func setupRouter(repo repository.DeviceRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := handlers.NewDeviceHandler(repo)
	r.GET("/devices", handler.GetDevices)
	r.POST("/devices", handler.CreateDevice)
	r.DELETE("/devices/:id", handler.DeleteDevice)
	return r
}

func TestGetDevices(t *testing.T) {
	mockRepo := &MockRepo{
		devices: []models.Device{
			{ID: "test-1", Name: "Test Device"},
		},
	}
	router := setupRouter(mockRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res []models.Device
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Len(t, res, 1)
	assert.Equal(t, "test-1", res[0].ID)
}

func TestCreateDevice(t *testing.T) {
	mockRepo := &MockRepo{}
	router := setupRouter(mockRepo)

	newDevice := models.Device{ID: "test-2", Name: "New Device"}
	body, _ := json.Marshal(newDevice)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Len(t, mockRepo.devices, 1)
	assert.Equal(t, "online", mockRepo.devices[0].Status)
}
