package api

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AppPort string `default:"8080" envconfig:"APP_PORT"`
}

func NewConfig(envPrefix string) (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
