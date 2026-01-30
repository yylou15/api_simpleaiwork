package main

import (
	"crypto/tls"
	"net/http"
	"net/url"
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
	sredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
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

	// Initialize Service and Handler
	userHandler := handler.NewUserHandler(service.NewUserService())
	paddleHandler := handler.NewPaddleHandler(service.NewUserService())

	// Public Routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/send-code", userHandler.SendVerificationCode)
		authGroup.POST("/login", userHandler.Login)
	}

	// Webhooks
	r.POST("/webhooks/paddle", paddleHandler.HandleWebhook)

	// Protected Routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", userHandler.GetUser)
		// Register can be protected or public depending on requirement.
		// Assuming public for now as it's often part of sign-up, but user said "homepage no need login", others need.
		// If Register is "admin create user", it should be protected. If "user sign up", public.
		// Given Login handles sign-up (FindOrCreate), POST /users might be redundant or for admin.
		// Let's keep POST /users public for now or put it in protected if it's strictly for authenticated users.
		// I'll leave it public for now to avoid breaking existing tests/usage, but technically usually sign-up is public.
	}
	r.POST("/users", userHandler.Register) // Kept public

	// Run the server on port 8080
	r.Run(":8080")
}

func buildRedisStore() sessions.Store {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("REDIS_URL environment variable is not set")
	}
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "secret123"
	}

	u, err := url.Parse(redisURL)
	if err != nil {
		panic(err)
	}
	password, _ := u.User.Password()

	useTLS := strings.HasPrefix(redisURL, "rediss://")
	pool := &redigo.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redigo.Conn, error) {
			options := []redigo.DialOption{}
			if password != "" {
				options = append(options, redigo.DialPassword(password))
			}
			if useTLS {
				options = append(options, redigo.DialTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}))
			}
			return redigo.Dial("tcp", u.Host, options...)
		},
		TestOnBorrow: func(conn redigo.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	store, err := sredis.NewStoreWithPool(pool, []byte(secret))
	if err != nil {
		panic(err)
	}

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
