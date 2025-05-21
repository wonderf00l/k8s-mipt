package logger

type Config struct {
	Level string `yaml:"level"`
}

func (c *Config) SetDefaults() {
	const (
		defaultLevel = "info"
	)

	if c.Level == "" {
		c.Level = defaultLevel
	}
}
