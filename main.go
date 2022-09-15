package main

import (
	"goshorten/internal/handler"
	"goshorten/internal/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	redis.Connect()
	defer redis.Disconnect()
	server := gin.Default()
	server.POST("/shorten", handler.Shorten)
  server.Run()
}