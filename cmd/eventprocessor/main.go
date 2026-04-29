package main

import (
	"log"
	"context"
	"strings"
	"pismo-take-home/config"
	"pismo-take-home/internal/consumer"
	"pismo-take-home/internal/eventprocessor"
	"pismo-take-home/internal/repository"
	"pismo-take-home/internal/validator"
	"pismo-take-home/schemas"
)

func main() {
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
	schemaLoader, error := schemas.NewSchemaLoader("schemas")

	if error != nil {
		log.Fatal(error)
	}

	payloadValidator := validator.NewPayloadValidator(schemaLoader)
	processor := eventprocessor.New(payloadValidator, repo)
	brokers := strings.Split(config.KafkaBrokers, ",")
	kafkaConsumer := consumer.NewKafkaConsumer(config.KafkaTopic, config.KafkaGroupID, brokers)
	ctx := context.Background()
	log.Println("starting event-processor")

	for {
		message, error := kafkaConsumer.ConsumeMessage(ctx)

		if error != nil {
			log.Println("Error fetching message:", error)
			continue
		}

		processingError := processor.ProcessEvent(ctx, message.ByteValue)

		if processingError != nil {
			log.Println("Error prossessing event:", processingError)
			continue
		}

		commitError := kafkaConsumer.CommitMessage(ctx,message)

		if commitError != nil {
			log.Println("Error committing Event:", commitError)
			continue
		}

		log.Println("event processed!")
	}
}