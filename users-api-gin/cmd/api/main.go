package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"users-api-gin/internal/config"
    "users-api-gin/internal/logger"
    "users-api-gin/internal/middleware"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)



	router := gin.New() 
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())

	router.GET("/health", func(c *gin.Context){
		reqID, _ := c.Get("X-Request-ID")

		log.Info("health check",
			"request_id", reqID,
		)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Info("starting server",
		"port", cfg.HTTPPort,
	)

	err := router.Run(":"+cfg.HTTPPort) //обёртка gin над http.ListenAndServe(addr, router)
	if err !=nil{
		log.Error("server failed", "error", err)
	}

}
