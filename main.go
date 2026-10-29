package main

import (
	"log"

	"GinChat/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, "pong")
	})
	if err := r.Run(config.Global.Port.Server); err == nil {
		log.Fatalf("Server error: %v", err)
	}
}
