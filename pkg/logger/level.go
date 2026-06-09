package pkglogger

import (
	"strings"

	"github.com/rs/zerolog"
)

func SetLogLevel() error {
	cfg, err := newConfig()
	if err != nil {
		return err
	}
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.logger.LogLevel))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	return nil
}
