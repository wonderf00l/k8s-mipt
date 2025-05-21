package http

type Config struct {
	SrvPort uint16 `yaml:"srvPort"`

	HandlerCfg HandlerConfig `yaml:"handler"`
}

func (c *Config) SetDefaults() {
	const (
		defaultPort     = 8080
	)

	if c.SrvPort == 0 {
		c.SrvPort = defaultPort
	}

	c.HandlerCfg.SetDefaults()
}
