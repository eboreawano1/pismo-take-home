package eventprocessor


import (
	"fmt"
	"context"
	"encoding/json"
	"pismo-take-home/internal/event"
	"pismo-take-home/internal/validator"
)

func New(dataStore DataStore) *EventProcessor {
	return &EventProcessor{ dataStore: dataStore }
}

func (eventProcessor *EventProcessor) ProcessEvent(ctx context.Context, eventBytes []byte) error {
	var event event.Event
	unmarshalError := json.Unmarshal(eventBytes, &event)

	if unmarshalError != nil {
		return fmt.Errorf("Error while unmarshalling event: %w", unmarshalError)
	}

	validationErrors := validator.Validate(event)
	
	if len(validationErrors) > 0 {
		return fmt.Errorf("Error: invalid event: %v", validationErrors)
	}

	saveError := eventProcessor.dataStore.Save(ctx, event)

	if saveError != nil {
		return fmt.Errorf("Error while saving event: %w", saveError)
	}

	return nil
}

type DataStore interface {
	Save(ctx context.Context, event event.Event) error
}

type EventProcessor struct {
	dataStore DataStore
}
