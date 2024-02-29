package config

import (
	"fmt"

	"github.com/protomem/gotube/pkg/env"
)

type Config struct {
	basePrefix string
	cache      map[string]any
}

func New() *Config {
	return &Config{
		basePrefix: "APP",
		cache:      make(map[string]any),
	}
}

func (c *Config) fmtPrefix(prefix string) string {
	return fmt.Sprintf("%s_%s_", c.basePrefix, prefix)
}

type configParser[T any] struct {
	cache map[string]any
}

func newConfigParser[T any](cache map[string]any) *configParser[T] {
	return &configParser[T]{
		cache: cache,
	}
}

func (cp *configParser[T]) parse(prefix string) (T, error) {
	var conf T
	if conf, ok := cp.cache[prefix]; ok {
		if conf, ok := conf.(T); ok {
			return conf, nil
		}
	}

	if err := env.ParseWithPrefix(&conf, prefix); err != nil {
		return conf, err
	}

	cp.cache[prefix] = conf
	return conf, nil
}
