package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTP struct {
		Addr string
	}

	Log struct {
		Level string
	}
}

func New() (Config, error) {
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

	return conf, errs
}
