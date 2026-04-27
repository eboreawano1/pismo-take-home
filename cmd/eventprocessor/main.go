package main

import (
	"log"
	"pismo-take-home/config"
)

func main() {
	log.Println("starting event-processor")
	config, error := config.Load()

	if(error != nil) {
		log.Fatal(error)
	}

	log.Println("configured database: ", config.DatabaseURL != "")
}