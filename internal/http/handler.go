package http

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	statusMsg = `{"status": "ok"}`
)

type (
	Handler struct {
		helloMsg string

		logDest *os.File

		mu sync.RWMutex
	}

	HandlerConfig struct {
		HelloMsg string `yaml:"helloMsg"`

		LogDir      string `yaml:"logDir"`
		LogFilename string `yaml:"logFilename"`
	}

	logMsg struct {
		Msg string `json:"message"`
	}
)

func (hc *HandlerConfig) SetDefaults() {
	const (
		defaultLogDir      = "app/logs/"
		defaultLogFilename = "app.log"

		defaultHelloMsg = "Welcome to the custom app"
	)

	if hc.LogDir == "" {
		hc.LogDir = defaultLogDir
	}
	if hc.LogFilename == "" {
		hc.LogFilename = defaultLogFilename
	}

	if hc.HelloMsg == "" {
		hc.HelloMsg = defaultHelloMsg
	}
}

func newHandler(cfg *HandlerConfig) (*Handler, error) {
	err := os.MkdirAll(cfg.LogDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("создание директории для логов: %w", err)
	}

	logFile, err := os.OpenFile(cfg.LogDir+cfg.LogFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("открытие лог-файла: %w", err)
	}

	return &Handler{
		helloMsg: cfg.HelloMsg,
		logDest:  logFile,
	}, nil
}

func (h *Handler) SayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h.helloMsg)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, statusMsg)
}

func (h *Handler) WriteLog(w http.ResponseWriter, r *http.Request) {
	var msg logMsg

	logBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Ошибка чтения тела запроса: %s", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err = json.Unmarshal(logBytes, &msg); err != nil {
		log.Error().Msgf("Ошибка десериализации тела сообщения для логирования: %s", err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	log.Info().Msgf("Получено сообщение: %s", msg)

	logBytes = append(logBytes, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, err = h.logDest.Write(logBytes); err != nil {
		log.Error().Msgf("Ошибка записи сообщения в лог-файл: %s", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetLog(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := h.getFileContent()
	if err != nil {
		log.Error().Msgf("Получение содержимого лог-файла: %s", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err = fmt.Fprint(w, fileBytes); err != nil {
		log.Error().Msgf("Ошибка записи содержимого файла в качестве тела ответа: %s", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *Handler) getFileContent() ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if _, err := h.logDest.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("установки file-pointer в начало файла: %w", err)
	}

	fileBytes, err := io.ReadAll(h.logDest)
	if err != nil {
		return nil, fmt.Errorf("чтение содержимого файла: %w", err)
	}

	return fileBytes, nil
}

func (h *Handler) Close() error {
	if err := h.logDest.Close(); err != nil {
		return fmt.Errorf("закрытие лог-файла: %w", err)
	}

	return nil
}

func (m *logMsg) String() string {
	return m.Msg
}
