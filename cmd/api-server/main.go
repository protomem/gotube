package main

import (
	"log"
	"os"

	application "github.com/protomem/gotube/internal/app"
	"github.com/protomem/gotube/internal/config"
)

func main() {
	var err error

	conf, err := config.New()
	if err != nil {
		log.Printf("config error: %v\n", err)
		os.Exit(1)
	}

	app, err := application.New(conf)
	if err != nil {
		log.Printf("application error: %v\n", err)
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		log.Printf("application error: %v\n", err)
		os.Exit(1)
	}
}
