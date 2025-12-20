package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"users-api-gin/internal/config"
    "users-api-gin/internal/logger"
    "users-api-gin/internal/middleware"
	"users-api-gin/internal/repository/postgres"
	"users-api-gin/internal/service"
	"users-api-gin/internal/handler"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)

	db, err := postgres.New(cfg.DBDSN)
	if err != nil {
		log.Error("failed to connect to db", "error", err)
		return
	}

	userRepo := postgres.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())

	// health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 🚨 ВАЖНО: USERS ROUTES
	router.POST("/users", userHandler.Create)
	router.GET("/users", userHandler.List)
	router.GET("/users/:id", userHandler.GetByID)

	log.Info("starting server", "port", cfg.HTTPPort)

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Error("server failed", "error", err)
	}
}


