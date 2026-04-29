package eventprocessor


import (
	"fmt"
	"context"
	"encoding/json"
	"pismo-take-home/internal/event"
	"pismo-take-home/internal/validator"
	"pismo-take-home/internal/router"
)

func New(payloadValidator *validator.PayloadValidator, dataStore DataStore) *EventProcessor {
	return &EventProcessor{ 
		dataStore: dataStore,
		payloadValidator: payloadValidator,
	}
}

func (eventProcessor *EventProcessor) ProcessEvent(ctx context.Context, eventBytes []byte) error {
	var e event.Event
	unmarshalError := json.Unmarshal(eventBytes, &e)

	if unmarshalError != nil {
		return fmt.Errorf("Error while unmarshalling event: %w", unmarshalError)
	}

	validationErrors := validator.ValidateEvent(e)
	validationErrors = append(validationErrors, eventProcessor.payloadValidator.ValidatePayload(e)...)
	
	if len(validationErrors) > 0 {
		saveError := eventProcessor.dataStore.Save(ctx, e, event.InvalidEventStatus, "", validationErrors)
		
		if saveError != nil {
			return fmt.Errorf("Error while saving invalid event: %w", saveError)
		}

		return nil
	}

	deliveryTarget := router.RouteEvent(e)

	if deliveryTarget == "" {
		processingError := eventProcessor.dataStore.Save(ctx, e, event.ProcessingErrorStatus, "", nil)
		
		if  processingError != nil {
			return fmt.Errorf("Error while processing unrouted event: %w", processingError)
		}

		return nil
	}

	saveError := eventProcessor.dataStore.Save(ctx, e, event.ReadyToDeliverStatus, deliveryTarget, validationErrors)

	if saveError != nil {
		return fmt.Errorf("Error while saving event: %w", saveError)
	}

	return nil
}

type DataStore interface {
	Save(
		ctx context.Context, 
		event event.Event,
		status string,
		deliveryTarget string,
		validationErrors []string,
	) error
}

type EventProcessor struct {
	dataStore DataStore
	payloadValidator *validator.PayloadValidator
}
