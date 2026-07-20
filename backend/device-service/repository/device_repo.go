package repository

import (
	"device-service/models"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetAll() ([]models.Device, error)
	Create(device *models.Device) error
	Delete(id string) error
	Count() (int64, error)
	UpdateStatus(id string, status string) error
}

type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) GetAll() ([]models.Device, error) {
	var devices []models.Device
	err := r.db.Find(&devices).Error
	return devices, err
}

func (r *deviceRepo) Create(device *models.Device) error {
	return r.db.Create(device).Error
}

func (r *deviceRepo) Delete(id string) error {
	return r.db.Delete(&models.Device{}, "id = ?", id).Error
}

func (r *deviceRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Device{}).Count(&count).Error
	return count, err
}

func (r *deviceRepo) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Device{}).Where("id = ?", id).Update("status", status).Error
}
