package config

import (
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var Global AppConfig

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	MaxConns int    `mapstructure:"max_conns"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Timeout struct {
	DelayHeartbeat   int `mapstructure:"delay_heartbeat"`
	HeartbeatHz      int `mapstructure:"heartbeat_hz"`
	HeartbeatMaxTime int `mapstructure:"heartbeat_max_time"`
	RedisOnlineTime  int `mapstructure:"redis_online_time"`
}

type Port struct {
	Server string `mapstructure:"server"`
}

type Log struct {
	Output     string `mapstructure:"output"`
	File       string `mapstructure:"file"`
	FormatJSON bool   `mapstructure:"format_json"`
}

type Jwt struct {
	AccessTokenSecret string `mapstructure:"access_token_secret"`
	AccessExpireTime  int    `mapstructure:"access_expire_time"`
}

type AppConfig struct {
	Postgres Postgres      `mapstructure:"postgres"`
	Redis    Redis         `mapstructure:"redis"`
	Timeout  Timeout       `mapstructure:"timeout"`
	Port     Port          `mapstructure:"port"`
	Log      Log           `mapstructure:"log"`
	DB       *gorm.DB      `mapstructure:"-"`
	Red      *redis.Client `mapstructure:"-"`
	Jwt      Jwt           `mapstructure:"jwt"`
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to Read Config:", err)
	}

	if err := viper.Unmarshal(&Global); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
}
