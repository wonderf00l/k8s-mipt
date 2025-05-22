package main

import (
	"context"
	"errors"
	"fmt"
	"k8s-mipt/internal/config"
	"k8s-mipt/internal/http"
	"k8s-mipt/internal/logger"
	_ "k8s-mipt/internal/metrics"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	cfgFileName = "config/app.yaml"
)

func run(ctx context.Context) error {
	cfg, err := config.New(cfgFileName)
	if err != nil {
		return fmt.Errorf("получение конфигурации приложения: %s", err)
	}

	if err = logger.SetupAppLogger(&cfg.Logging); err != nil {
		return fmt.Errorf("конфигурирование логгера: %s", err)
	}

	srv, err := http.NewServer(&cfg.Http)
	if err != nil {
		return fmt.Errorf("создание http сервера: %s", err)
	}

	if err = srv.Start(); err != nil {
		return fmt.Errorf("старт http сервера: %s", err)
	}

	<-ctx.Done()

	log.Info().Msgf("Старт gracefull shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errs := []error{}

	if err = srv.Close(shutdownCtx); err != nil {
		errs = append(errs, fmt.Errorf("закрытие сервера: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("исполнение gracefull shutdown: %w", errors.Join(errs...))
	}

	log.Debug().Msg("Остановка приложения завершена успешно")

	return nil
}
