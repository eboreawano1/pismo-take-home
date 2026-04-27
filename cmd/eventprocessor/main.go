package main

import (
	"log"
	"context"
	"pismo-take-home/config"
	"pismo-take-home/internal/eventprocessor"
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
	processor := eventprocessor.New(repo)

	eventBytes := []byte(`{
		"event_id": "1",
		"event_type": "TEST_EVENT",
		"tenant_id": "tenant-1",
		"producer": "MANUAL",
		"event_time": "2026-04-27T00:00:00Z",
		"schema_version": "1",
		"payload": {}
	}`)

	processingError := processor.ProcessEvent(context.Background(), eventBytes)

	if processingError != nil {
		log.Fatal(processingError)
	}

	log.Println("event processed!")
}