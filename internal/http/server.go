package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Srv struct {
	http.Server

	handler *Handler

	done chan struct{}
}

func (s *Srv) Start() error {
	var (
		err     error
		serveCh = make(chan error, 1)
	)

	go func() {
		close(s.done)

		log.Info().Msgf("Запуск http сервера на %s", s.Server.Addr)

		err = s.ListenAndServe()
		if err != nil {
			serveCh <- err
		}
	}()

	const httpSrvStartTime = 200 * time.Millisecond

	select {
	case <-serveCh:
		return fmt.Errorf("исполнение listen and serve: %w", err)
	case <-time.After(httpSrvStartTime):
	}

	return nil
}

func (s *Srv) Close(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown http сервера: %w", err)
	}

	select {
	case <-s.done:
	case <-ctx.Done():
		return fmt.Errorf("ожидание listener потока: %w", ctx.Err())
	}

	if err = s.handler.Close(); err != nil {
		return fmt.Errorf("закрытие handler: %w", err)
	}

	log.Info().Msg("Закрытие http сервера успешно завершено")

	return nil
}

func NewServer(cfg *Config) (*Srv, error) {
	handler, err := newHandler(&cfg.HandlerCfg)
	if err != nil {
		return nil, fmt.Errorf("создание http handler: %w", err)
	}

	return &Srv{
		Server: http.Server{
			Addr:    ":" + strconv.Itoa(int(cfg.SrvPort)),
			Handler: newRouter(handler),
		},
		done:    make(chan struct{}),
		handler: handler,
	}, nil
}

func newRouter(handler *Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", handler.SayHello)
	router.Get("/status", handler.Status)
	router.Get("/logs", handler.GetLog)
	router.Post("/log", handler.WriteLog)

	return router
}
