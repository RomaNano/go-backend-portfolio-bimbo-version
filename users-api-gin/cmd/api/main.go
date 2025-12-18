package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"users-api-gin/internal/config"
)

func main() {
	cfg := config.Load()

	router := gin.New() 
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("starting server on :%s", cfg.HTTPPort)
	err := router.Run(":"+cfg.HTTPPort)
	if err !=nil{
		log.Fatal(err)
	}

}
