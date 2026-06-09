package pkglogger

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:""`
}

func newConfig() (*config, error) {
	_ = godotenv.Load()

	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
