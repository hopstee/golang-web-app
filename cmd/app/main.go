package main

import (
	"log"
	"mobile-backend-boilerplate/internal/app"
)

func main() {
	if err := Run(); err != nil {
		log.Fatalf("program finished with error: %v", err)
	}
}

func Run() error {
	app, err := app.Init()
	if err != nil {
		return err
	}
	defer app.Close()

	return app.Run()
}
