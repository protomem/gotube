package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	HTTP struct {
		Addr string
	}

	Log struct {
		Level string
	}

	Auth struct {
		Secret string
	}
}

func New() (Config, error) {
	const op = "config.New"
	var (
		conf Config
		ok   bool
		errs error
	)

	conf.HTTP.Addr, ok = os.LookupEnv("HTTP__ADDR")
	if !ok {
		errs = errors.Join(errs, errors.New("HTTP__ADDR is not set"))
	}

	conf.Log.Level, ok = os.LookupEnv("LOG__LEVEL")
	if !ok {
		errs = errors.Join(errs, errors.New("LOG__LEVEL is not set"))
	}

	conf.Auth.Secret, ok = os.LookupEnv("AUTH__SECRET")
	if !ok {
		errs = errors.Join(errs, errors.New("AUTH__SECRET is not set"))
	}

	if errs != nil {
		return Config{}, fmt.Errorf("%s: %w", op, errs)
	}

	return conf, nil
}
