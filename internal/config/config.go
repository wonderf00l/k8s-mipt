package config

import (
	"fmt"
	"k8s-mipt/internal/http"
	"k8s-mipt/internal/logger"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Http    http.Config   `yaml:"http"`
		Logging logger.Config `yaml:"logging"`
	}
)

func New(cfgFileName string) (*Config, error) {
	cfgBytes, err := os.ReadFile(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("чтение файла конфигурации: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		return nil, fmt.Errorf("десериализация конфигурации: %w", err)
	}

	cfg.setDefaults()

	return &cfg, nil
}

func (c *Config) setDefaults() {
	c.Http.SetDefaults()
	c.Logging.SetDefaults()
}
