package env

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func Load(filenames ...string) error {
	if len(filenames) == 0 {
		return nil
	}

	if err := godotenv.Load(filenames...); err != nil {
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
