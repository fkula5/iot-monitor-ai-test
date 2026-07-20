package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var jwtKey = []byte("super_secret_key_123") // In prod, load from ENV

type RoleDefinition struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	Name            string `json:"name" gorm:"unique"`
	CanWriteDevices bool   `json:"canWriteDevices"`
	CanWriteRules   bool   `json:"canWriteRules"`
	CanManageUsers  bool   `json:"canManageUsers"`
}

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"`
	Role     string `json:"role"` // This corresponds to RoleDefinition.Name
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username        string `json:"username"`
	Role            string `json:"role"`
	CanWriteDevices bool   `json:"canWriteDevices"`
	CanWriteRules   bool   `json:"canWriteRules"`
	CanManageUsers  bool   `json:"canManageUsers"`
	jwt.RegisteredClaims
}

var db *gorm.DB

func initDB() {
	var err error
	
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" { dbHost = "localhost" }
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" { dbPort = "5432" }
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" { dbUser = "admin" }
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" { dbPass = "adminpassword" }
	dbName := os.Getenv("DB_NAME")
	if dbName == "" { dbName = "iot_db" }

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Warsaw", dbHost, dbUser, dbPass, dbName, dbPort)
	
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Println("waiting for postgres...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect database after retries: ", err)
	}
	db.AutoMigrate(&RoleDefinition{}, &User{})

	// Seed roles
	var roleCount int64
	db.Model(&RoleDefinition{}).Count(&roleCount)
	if roleCount == 0 {
		db.Create(&RoleDefinition{Name: "Admin", CanWriteDevices: true, CanWriteRules: true, CanManageUsers: true})
		db.Create(&RoleDefinition{Name: "Viewer", CanWriteDevices: false, CanWriteRules: false, CanManageUsers: false})
	}

	// Create demo user
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&User{Username: "admin", Password: string(hashed), Role: "Admin"})
	} else {
		// Migrate existing admin to Admin role
		db.Model(&User{}).Where("username = ?", "admin").Update("role", "Admin")
	}
}

func loginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var roleDef RoleDefinition
	db.Where("name = ?", user.Role).First(&roleDef)

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username:        creds.Username,
		Role:            user.Role,
		CanWriteDevices: roleDef.CanWriteDevices,
		CanWriteRules:   roleDef.CanWriteRules,
		CanManageUsers:  roleDef.CanManageUsers,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func registerHandler(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	user := User{Username: creds.Username, Password: string(hashed), Role: "Viewer"}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username taken"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func getUsersHandler(c *gin.Context) {
	var users []User
	if err := db.Select("id", "username", "role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func updateRoleHandler(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Role string `json:"role"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Validate if role exists
	var roleDef RoleDefinition
	if err := db.Where("name = ?", payload.Role).First(&roleDef).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role does not exist"})
		return
	}

	if err := db.Model(&User{}).Where("id = ?", id).Update("role", payload.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func deleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// Role Management Handlers
func getRolesHandler(c *gin.Context) {
	var roles []RoleDefinition
	if err := db.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func createRoleHandler(c *gin.Context) {
	var payload RoleDefinition
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := db.Create(&payload).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to create role, might already exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role created", "role": payload})
}

func deleteRoleHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "Admin" || name == "Viewer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete system default roles"})
		return
	}
	if err := db.Where("name = ?", name).Delete(&RoleDefinition{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)
	r.GET("/users", getUsersHandler)
	r.PUT("/users/:id/role", updateRoleHandler)
	r.DELETE("/users/:id", deleteUserHandler)

	r.GET("/roles", getRolesHandler)
	r.POST("/roles", createRoleHandler)
	r.DELETE("/roles/:name", deleteRoleHandler)
	
	return r
}

func main() {
	initDB()
	r := setupRouter()

	log.Println("Auth Service running on port 8083")
	r.Run(":8083")
}
