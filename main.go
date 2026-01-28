package main

import (
	"net/http"
	"strings"

	"api/cert"
	"api/database"
	"api/internal/dal/query"
	"api/internal/handler"
	"api/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	cert.Init()
	// Load environment variables
	godotenv.Load()

	// Initialize Database
	database.Connect()

	// Initialize GORM Gen Query
	query.SetDefault(database.DB)

	// Initialize Gin engine
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		// Allow root domain (http and https)
		if origin == "https://simpleaiwork.com" || origin == "http://simpleaiwork.com" {
			return true
		}
		// Allow subdomains
		if strings.HasSuffix(origin, ".simpleaiwork.com") && (strings.HasPrefix(origin, "http://") || strings.HasPrefix(origin, "https://")) {
			return true
		}
		return false
	}
	r.Use(cors.New(config))

	// Define a simple GET route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong2",
		})
	})

	// Initialize Service and Handler
	userHandler := handler.NewUserHandler(service.NewUserService())
	// User Routes
	r.POST("/users", userHandler.Register)
	r.GET("/users", userHandler.GetUser)

	// Run the server on port 8080
	r.Run(":8080")
}
