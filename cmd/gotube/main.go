package main

import (
	"flag"
	"log"
	"os"

	"github.com/protomem/gotube/internal/app"
	"github.com/protomem/gotube/internal/config"
)

var _confFile string

func init() {
	flag.StringVar(&_confFile, "conf", "", "path to config file")
}

func main() {
	var err error

	flag.Parse()
	if _confFile == "" {
		log.Println("error: no config file")
		os.Exit(1)
	}

	conf, err := config.New(_confFile)
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	app, err := app.New(conf)
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
