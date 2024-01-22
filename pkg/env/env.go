package env

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func Load(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("env.Load: %w", err)
	}
	return nil
}

func Parse(v any) error {
	if err := env.Parse(v); err != nil {
		return fmt.Errorf("env.Parse: %w", err)
	}
	return nil
}
