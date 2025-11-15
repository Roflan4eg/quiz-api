package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
