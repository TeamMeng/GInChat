package initial

import (
	"context"
	"log"

	"GinChat/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	config.Global.Red = redis.NewClient(&redis.Options{
		Addr:     config.Global.Redis.Addr,
		Password: config.Global.Redis.Password,
		DB:       config.Global.Redis.DB,
	})

	if _, err := config.Global.Red.Ping(context.Background()).Result(); err != nil {
		log.Fatal("Failed to connect Redis:", err)
	}
}
