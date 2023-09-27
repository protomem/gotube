package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"addr" env-required:"true"`
	} `yaml:"http"`

	Auth struct {
		Secret string `yaml:"secret" env-required:"true"`
	} `yaml:"auth"`

	Log struct {
		Level string `yaml:"level" env-required:"true"`
	} `yaml:"log"`

	Postgres struct {
		Connect string `yaml:"connect" env-required:"true"`
	} `yaml:"postgres"`

	Mongo struct {
		Connect string `yaml:"connect" env-required:"true"`
	} `yaml:"mongo"`

	S3 struct {
		Addr      string `yaml:"addr" env-required:"true"`
		AccessKey string `yaml:"access_key" env-required:"true"`
		SecretKey string `yaml:"secret_key" env-required:"true"`
	} `yaml:"s3"`

	Redis struct {
		Connect string `yaml:"connect" env-required:"true"`
	} `yaml:"redis"`
}

func New(filename string) (Config, error) {
	var conf Config

	err := cleanenv.ReadConfig(filename, &conf)
	if err != nil {
		return Config{}, fmt.Errorf("config.New(%s): %w", filename, err)
	}

	return conf, nil
}
