package api

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App AppConfig
}
type AppConfig struct {
	Host        string `envconfig:"APP_HOST" default:"localhost"`
	Port        string `envconfig:"APP_PORT" default:"8081"`
	ServiceName string `envconfig:"APP_SERVICE_NAME" default:"golang-backend-k03-with-echo-v4"`
	InstanceID  string `envconfig:"APP_INSTANCE_ID" default:""`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {

		// Nếu parse lỗi (type mismatch, thiếu field bắt buộc...)
		return nil, err
	}
	return &cfg, nil
}
