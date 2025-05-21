package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupAppLogger(cfg *Config) error {
	lvl, err := stringToZerologLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("конвертация уровня логирования %s: %w", cfg.Level, err)
	}

	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger().
		Level(lvl)

	return nil
}

func stringToZerologLevel(level string) (zerolog.Level, error) {
	switch level {
	case "info", "":
		return zerolog.InfoLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "trace":
		return zerolog.TraceLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	default:
		return 0, fmt.Errorf("передано некорректное строковое значение уровня логирования %s", level)
	}
}
