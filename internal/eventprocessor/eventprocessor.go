package eventprocessor


import (
	"fmt"
	"context"
	"encoding/json"
	"pismo-take-home/internal/event"
	"pismo-take-home/internal/validator"
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
		saveError := eventProcessor.dataStore.Save(ctx, e, event.InvalidEventStatus, validationErrors)
		
		if saveError != nil {
			return fmt.Errorf("Error while saving event: %w", saveError)
		}

		return nil
	}

	saveError := eventProcessor.dataStore.Save(ctx, e, event.ReadyToDeliverStatus, validationErrors)

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
		validationErrors []string,
	) error
}

type EventProcessor struct {
	dataStore DataStore
	payloadValidator *validator.PayloadValidator
}
