package config

import (
	"fmt"

	"github.com/protomem/gotube/pkg/env"
)

type HTTP struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"8080"`
}

type Log struct {
	Level string `env:"LEVEL" envDefault:"debug"`
}

type Database struct {
	DSN string `env:"DSN,notEmpty"`
}

type Flash struct {
	DSN string `env:"DSN" envDefault:"localhost:6379/0"`
}

type Config struct {
	HTTP     `envPrefix:"HTTP_"`
	Log      `envPrefix:"LOG_"`
	Database `envPrefix:"DB_"`
	Flash    `envPrefix:"FLASH_"`
}

func New() (Config, error) {
	var conf Config
	if err := env.Parse(&conf); err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}
	return conf, nil
}
