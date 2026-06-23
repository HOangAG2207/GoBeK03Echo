package api

import (
	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// App  AppConfig
	Host        string `envconfig:"APP_HOST" default:"localhost"`
	Port        string `envconfig:"APP_PORT" default:"8080"`
	ServiceName string `envconfig:"APP_SERVICE_NAME" default:"golang-backend-k03-with-echo-v4"`
	InstanceID  string `envconfig:"APP_INSTANCE_ID" default:""`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {

		// Nếu parse lỗi (type mismatch, thiếu field bắt buộc...)
		return nil, err
	}
	if cfg.InstanceID == "" {
		cfg.InstanceID = helpers.UuidGenerator()
	}
	return &cfg, nil
}
