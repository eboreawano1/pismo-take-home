package eventprocessor


import (
	"fmt"
	"context"
	"encoding/json"
	"pismo-take-home/internal/event"
	"pismo-take-home/internal/repository"
)

func New(repository *repository.Repository) *EventProcessor {
	return &EventProcessor{ repository: repository }
}

func (eventProcessor *EventProcessor) ProcessEvent(ctx context.Context, eventBytes []byte) error {
	var event event.Event
	unmarshalError := json.Unmarshal(eventBytes, &event)

	if unmarshalError != nil {
		return fmt.Errorf("Error while unmarshalling event: %w", unmarshalError)
	}

	saveError := eventProcessor.repository.Save(ctx, event)

	if saveError != nil {
		return fmt.Errorf("Error while saving event: %w", saveError)
	}

	return nil
}

type EventProcessor struct {
	repository *repository.Repository
}
