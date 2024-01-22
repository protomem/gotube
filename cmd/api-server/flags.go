package main

import (
	"flag"

	"github.com/protomem/gotube/pkg/env"
)

func init() {
	confFile := flag.String("conf", "", "path to config file")

	flag.Parse()

	if *confFile != "" {
		if err := env.Load(*confFile); err != nil {
			panic(err)
		}
	}
}
