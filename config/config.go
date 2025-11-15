package config

import (
	"fmt"
	"time"
)

type HTTPConfig struct {
	Port         string        `env:"PORT" envDefault:"8080"`
	Host         string        `env:"HOST" envDefault:"0.0.0.0"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
}

func (c *HTTPConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type AppConfig struct {
	Name            string        `env:"NAME"`
	Environment     string        `env:"ENV" envDefault:"local"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	DBurl           string        `env:"DB_URL"`
}

type Config struct {
	App  AppConfig  `envPrefix:"APP_"`
	HTTP HTTPConfig `envPrefix:"HTTP_"`
}
