package main

import (
	"log"

	"GinChat/config"
	"GinChat/initial"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	initial.InitRedis()
	initial.InitLogger()
	if err := initial.InitPostgres(); err != nil {
		log.Fatalf("Failed to connect postgres: %v", err)
	}

	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, "pong")
	})
	if err := r.Run(config.Global.Port.Server); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
