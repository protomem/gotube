package main

import (
	"log"
	"os"

	"github.com/protomem/gotube/internal/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Printf("App Failed: %v", err)
		os.Exit(1)
	}
}
