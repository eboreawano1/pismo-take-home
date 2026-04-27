package main

import (
	"log"
	"pismo-take-home/config"
	"pismo-take-home/internal/repository"
)

func main() {
	log.Println("starting event-processor")
	config, error := config.Load()

	if(error != nil) {
		log.Fatal(error)
	}

	log.Println("configured database: ", config.DatabaseURL != "")
	repo, error := repository.New(config.DatabaseURL)

	if error != nil {
		log.Fatal(error)
	}

	defer repo.Close()
	log.Println("database connected")
}