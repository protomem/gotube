package config

import (
	"errors"
	"fmt"

	"github.com/protomem/gotube/pkg/env"
)

const (
	Dev  Mode = "dev"
	Prod Mode = "prod"
)

type Mode string

func (m Mode) Validate() error {
	switch m {
	case Dev, Prod:
		return nil
	default:
		return errors.New("invalid mode")
	}
}

type Config struct {
	Mode Mode `env:"MODE" envDefault:"dev"`

	HTTP struct {
		Host string `env:"HOST" envDefault:"localhost"`
		Port int    `env:"PORT" envDefault:"8080"`
	} `envPrefix:"HTTP_"`

	Auth struct {
		Secret string `env:"SECRET,notEmpty,unset"`
	} `envPrefix:"AUTH_"`

	Log struct {
		Level string `env:"LEVEL" envDefault:"debug"`
	} `envPrefix:"LOG_"`

	Postgres struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     int    `env:"PORT" envDefault:"5432"`
		User     string `env:"USER,notEmpty"`
		Password string `env:"PASSWORD,notEmpty,unset"`
		Database string `env:"DATABASE" envDefault:"gotubedb"`
		Secure   bool   `env:"SECURE" envDefault:"false"`
	} `envPrefix:"POSTGRES_"`

	Mongo struct {
		URI string `env:"URI,notEmpty"`
	} `envPrefix:"MONGO_"`

	Redis struct {
		Addr string `env:"ADDR" envDefault:"localhost:6379"`
	} `envPrefix:"REDIS_"`

	S3 struct {
		Addr   string `env:"ADDR" envDefault:"localhost:9000"`
		Access string `env:"ACCESS,notEmpty"`
		Secret string `env:"SECRET,notEmpty,unset"`
		Secure bool   `env:"SECURE" envDefault:"false"`
	} `envPrefix:"S3_"`
}

func New() (Config, error) {
	const op = "config:Config.New"
	var conf Config

	if err := env.Parse(&conf); err != nil {
		return Config{}, fmt.Errorf("%s: %w", op, err)
	}

	return conf, nil
}
