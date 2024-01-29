package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/protomem/gotube/pkg/env"
)

type HTTP struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"8080"`
}

func (c HTTP) Addr() string { return fmt.Sprintf("%s:%s", c.Host, strconv.Itoa(c.Port)) }

type Log struct {
	Level string `env:"LEVEL" envDefault:"debug"`
}

type Auth struct {
	Secret          string        `env:"SECRET,notEmpty"`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"3h"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"24h"`
	Issuer          string        `env:"ISSUER" envDefault:"gotube"`
}

type Database struct {
	DSN string `env:"DSN,notEmpty"`
}

type Flash struct {
	DSN string `env:"DSN" envDefault:"localhost:6379/0"`
}

type Blob struct {
	Addr      string `env:"ADDR" envDefault:"localhost:9000"`
	AccessKey string `env:"ACCESS_KEY,notEmpty"`
	SecretKey string `env:"SECRET_KEY,notEmpty"`
	Secure    bool   `env:"SECURE" envDefault:"false"`
}

type Config struct {
	HTTP     `envPrefix:"HTTP_"`
	Log      `envPrefix:"LOG_"`
	Auth     `envPrefix:"AUTH_"`
	Database `envPrefix:"DB_"`
	Flash    `envPrefix:"FLASH_"`
	Blob     `envPrefix:"BLOB_"`
}

func New() (Config, error) {
	var conf Config
	if err := env.Parse(&conf); err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}
	return conf, nil
}
