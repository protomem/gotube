package main

import (
	application "github.com/protomem/gotube/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(application.Create()).Run()
}
