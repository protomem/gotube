package main

import (
	"flag"

	"github.com/protomem/gotube/pkg/env"
)

func init() {
	var confFile = flag.String("conf", "", "path to config file")

	flag.Parse()

	if *confFile != "" {
		_ = env.Load(*confFile)
	}
}
