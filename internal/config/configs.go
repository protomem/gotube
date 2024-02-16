package config

import (
	"fmt"
	"time"
)

type Server struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"8080"`
}

func (c *Config) Server() (Server, error) {
	prefix := "SERVER"
	conf, err := newConfigParser[Server](c.cache).parse(c.fmtPrefix(prefix))
	if err != nil {
		return conf, fmt.Errorf("config.%s: %w", prefix, err)
	}
	return conf, nil
}

type Log struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

func (c *Config) Log() (Log, error) {
	prefix := "LOG"
	conf, err := newConfigParser[Log](c.cache).parse(c.fmtPrefix(prefix))
	if err != nil {
		return conf, fmt.Errorf("config.%s: %w", prefix, err)
	}
	return conf, nil
}

type SQLiteDB struct {
	DSN string `env:"DSN" envDefault:":memory:"`
}

func (c *Config) SQLiteDB() (SQLiteDB, error) {
	prefix := "SQLITE"
	conf, err := newConfigParser[SQLiteDB](c.cache).parse(c.fmtPrefix(prefix))
	if err != nil {
		return conf, fmt.Errorf("config.%s: %w", prefix, err)
	}
	return conf, nil
}

type Auth struct {
	Secret         string        `env:"SECRET" envDefault:"secret"`
	AccessTokenTTL time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"72h"`
}

func (c *Config) Auth() (Auth, error) {
	prefix := "AUTH"
	conf, err := newConfigParser[Auth](c.cache).parse(c.fmtPrefix(prefix))
	if err != nil {
		return conf, fmt.Errorf("config.%s: %w", prefix, err)
	}
	return conf, nil
}
