package main

import (
	"log"
	"mobile-backend-boilerplate/internal/scripts"
)

func main() {
	if err := scripts.Execute(); err != nil {
		log.Fatal(err)
	}
}
