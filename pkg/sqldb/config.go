package pkgdb

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"admin"`
	Password string `envconfig:"DB_PASSWORD" default:""`
	DBName   string `envconfig:"DB_NAME" default:"bookmark"`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
	Timezone string `envconfig:"DB_TIMEZONE" default:"UTC"`
}

func newConfig(envprefix string) (*config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	var cfg config
	if err := envconfig.Process(envprefix, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
		cfg.Timezone,
	)
}
