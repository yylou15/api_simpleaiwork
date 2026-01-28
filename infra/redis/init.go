package redis

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is not set")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL:", err)
	}

	// DigitalOcean Managed Redis requires TLS
	// ParseURL should handle rediss:// scheme and enable TLS,
	// but sometimes we need to ensure TLS config is set correctly if using self-signed certs or specific constraints.
	// For managed services with valid public certs, default TLS config is usually fine.
	// If insecure skip verify is needed:
	// opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// For now, we trust the standard CA roots or the provider's setup.

	// Ensure TLS is enabled if scheme was rediss
	if opt.TLSConfig == nil && (redisURL[:8] == "rediss:/" || redisURL[:8] == "rediss:/") {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	Client = redis.NewClient(opt)

	_, err = Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Redis connection established successfully")
}
