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

	Postgres struct {
		Connect string
	}

	Mongo struct {
		URI string
	}

	Redis struct {
		Addr string
	}

	S3 struct {
		Addr string

		Keys struct {
			Access string
			Secret string
		}
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

	conf.Postgres.Connect, ok = os.LookupEnv("POSTGRES__CONNECT")
	if !ok {
		errs = errors.Join(errs, errors.New("POSTGRES__CONNECT is not set"))
	}

	conf.Mongo.URI, ok = os.LookupEnv("MONGO__URI")
	if !ok {
		errs = errors.Join(errs, errors.New("MONGO__URI is not set"))
	}

	conf.Redis.Addr, ok = os.LookupEnv("REDIS__ADDR")
	if !ok {
		errs = errors.Join(errs, errors.New("REDIS__ADDR is not set"))
	}

	conf.S3.Addr, ok = os.LookupEnv("S3__ADDR")
	if !ok {
		errs = errors.Join(errs, errors.New("S3__ADDR is not set"))
	}

	conf.S3.Keys.Access, ok = os.LookupEnv("S3__KEYS__ACCESS")
	if !ok {
		errs = errors.Join(errs, errors.New("S3__KEYS__ACCESS is not set"))
	}

	conf.S3.Keys.Secret, ok = os.LookupEnv("S3__KEYS__SECRET")
	if !ok {
		errs = errors.Join(errs, errors.New("S3__KEYS__SECRET is not set"))
	}

	if errs != nil {
		return Config{}, fmt.Errorf("%s: %w", op, errs)
	}

	return conf, nil
}
