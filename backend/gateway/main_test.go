package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(permissions map[string]bool) string {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"username": "testuser",
		"role":     "Viewer",
		"exp":      expirationTime.Unix(),
	}
	for k, v := range permissions {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func TestJWTMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/devices", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequirePermission_HasPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	token := generateTestToken(map[string]bool{"canWriteDevices": true})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/devices", nil) // /api/devices POST requires canWriteDevices
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Will return 500 because it attempts to proxy and device-service might be unreachable in this test,
	// but it should NOT return 401 or 403.
	assert.NotEqual(t, http.StatusForbidden, w.Code)
	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
}

func TestRequirePermission_MissingPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// Missing canWriteDevices
	token := generateTestToken(map[string]bool{"canWriteRules": true})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/devices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
