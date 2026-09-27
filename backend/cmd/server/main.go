package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/example/fullstack-assessment/backend/internal/handler"
	"github.com/example/fullstack-assessment/backend/internal/repository"
	"github.com/example/fullstack-assessment/backend/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

func main() {
	port := getenv("PORT", "8080")

	db, err := sql.Open(
		"mysql",
		getenv(
			"DB_DSN",
			"taskuser:taskpass@tcp(127.0.0.1:3306)/task_management?parseTime=true&charset=utf8mb4",
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("mysql: ", err)
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     getenv("REDIS_ADDR", "127.0.0.1:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       intEnv("REDIS_DB", 0),
	})

	if err := rc.Ping(context.Background()).Err(); err != nil {
		log.Printf("redis unavailable: %v (cache disabled)", err)
		rc = nil
	}

	ttl := time.Duration(intEnv("CACHE_TTL_SECONDS", 60)) * time.Second

	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo, rc, ttl)
	h := handler.NewTaskHandler(svc)

	r := gin.Default()

	// CORS middleware
	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	handler.Register(api, h)

	log.Printf("API running on :%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Izinkan frontend Expo Web
		if origin == "http://localhost:8081" ||
			origin == "http://localhost:8082" ||
			origin == "http://127.0.0.1:8081" ||
			origin == "http://127.0.0.1:8082" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func intEnv(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
