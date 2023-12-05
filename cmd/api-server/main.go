package main

import (
	"context"
	"flag"
	"fmt"

	application "github.com/protomem/gotube/internal/app"
	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/pkg/env"
)

var _confFile = flag.String("conf", "", "path to config file")

func init() {
	flag.Parse()
}

func main() {
	ctx := context.Background()

	if *_confFile != "" {
		if err := env.Load(*_confFile); err != nil {
			panic(fmt.Sprintf("failed load env from file %s: %v", *_confFile, err))
		}
	}

	conf, err := config.New()
	if err != nil {
		panic(fmt.Sprintf("failed init config: %v", err))
	}

	if err := application.New(conf).Run(ctx); err != nil {
		panic(fmt.Sprintf("failed run application: %v", err))
	}
}
