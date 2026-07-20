package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	db.AutoMigrate(&RoleDefinition{}, &User{})

	db.Create(&RoleDefinition{Name: "Admin", CanWriteDevices: true, CanWriteRules: true, CanManageUsers: true})
	db.Create(&RoleDefinition{Name: "Viewer", CanWriteDevices: false, CanWriteRules: false, CanManageUsers: false})
}

func TestRegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	router := setupRouter()

	// 1. Register
	creds := map[string]string{"username": "testuser", "password": "testpassword"}
	body, _ := json.Marshal(creds)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Login
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	
	var res map[string]string
	json.Unmarshal(w2.Body.Bytes(), &res)
	
	assert.NotEmpty(t, res["token"])
}

func TestDeleteAdminRoleReturnsError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/roles/Admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
