package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"api/biz/say_right/dal/query"
	"api/biz/say_right/handler"
	"api/biz/say_right/service"
	"api/cert"
	"api/database"
	"api/infra/mail"
	"api/infra/redis"
	"api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	cookieStore "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	cert.Init()
	// Load environment variables
	godotenv.Load()
	mail.CreateEmailClient()
	redis.Init()

	// Initialize Database
	database.Connect()

	// Initialize GORM Gen Query
	query.SetDefault(database.DB)

	// Initialize Gin engine
	r := gin.Default()

	store := buildRedisStore()
	r.Use(sessions.Sessions("mysession", store))

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.AllowOriginFunc = func(origin string) bool {
		// Allow localhost for development
		if strings.HasPrefix(origin, "http://localhost") {
			return true
		}
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

	// Register product routes
	registerProductRoutes(r)

	// Run the server on port 8080
	r.Run(":8080")
}

// registerProductRoutes registers routes for all products
func registerProductRoutes(r *gin.Engine) {
	// Initialize Paddle handler for global webhooks
	paddleHandler := handler.NewPaddleHandler(service.NewUserService())

	// Global webhooks
	r.POST("/webhooks/paddle", paddleHandler.HandleWebhook)

	// Register sayright routes
	registerSayRightRoutes(r)

	// Future products can be added here
	// registerCareerGameRoutes(r)
	// registerOtherProductRoutes(r)
}

// registerSayRightRoutes registers routes for sayright product
func registerSayRightRoutes(r *gin.Engine) {
	// Initialize Service and Handler
	userHandler := handler.NewUserHandler(service.NewUserService())
	templateHandler := handler.NewTemplateHandler(service.NewTemplateService())

	// Create product route group with prefix
	sayRightGroup := r.Group("/sayright")
	{
		// Public Routes
		authGroup := sayRightGroup.Group("/auth")
		{
			authGroup.POST("/send-code", userHandler.SendVerificationCode)
			authGroup.POST("/login", userHandler.Login)
		}

		// Protected Routes
		protected := sayRightGroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users", userHandler.GetUser)
			protected.GET("/templates", templateHandler.ListTemplates)
			protected.GET("/templates/:id", templateHandler.GetTemplateDetail)
		}

		// Public route
		sayRightGroup.POST("/users", userHandler.Register)
	}
}

func buildRedisStore() sessions.Store {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "secret123"
	}

	// Use cookie store for testing
	store := cookieStore.NewStore([]byte(secret))

	maxAge := 86400 * 7
	if v := os.Getenv("SESSION_MAX_AGE"); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			maxAge = int(parsed.Seconds())
		}
	}

	options := sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	}
	if domain := os.Getenv("SESSION_DOMAIN"); domain != "" {
		options.Domain = domain
	}
	if os.Getenv("SESSION_SECURE") == "true" {
		options.Secure = true
	}
	if strings.EqualFold(os.Getenv("SESSION_SAMESITE"), "none") {
		options.SameSite = http.SameSiteNoneMode
		options.Secure = true
	} else if strings.EqualFold(os.Getenv("SESSION_SAMESITE"), "strict") {
		options.SameSite = http.SameSiteStrictMode
	} else {
		options.SameSite = http.SameSiteLaxMode
	}

	store.Options(options)
	return store
}
