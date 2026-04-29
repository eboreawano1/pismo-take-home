package eventprocessor


import (
	"testing"
	"context"
	"errors"
	"pismo-take-home/internal/event"
	"pismo-take-home/schemas"
	"pismo-take-home/internal/validator"
)



func newTestProcessor(test *testing.T, dataStore DataStore) *EventProcessor {
	test.Helper()
	schemaLoader, error := schemas.NewSchemaLoader("../../schemas")

	if error != nil {
		test.Fatalf("Error creating schema loader: %v", error)
	}

	payloadValidator := validator.NewPayloadValidator(schemaLoader)

	return New(payloadValidator, dataStore)
}

func (dataStore *MockDataStore) Save(ctx context.Context, event event.Event) error {
	if dataStore.error != nil {
		return dataStore.error
	}

	dataStore.persistedEvents = append(dataStore.persistedEvents, event)

	return nil
}

func TestProcessEvent_ValidEvent_PersistsEvent(test *testing.T) {
	dataStore := &MockDataStore{}
	processor := newTestProcessor(test, dataStore)

	eventBytes := []byte(`{
		"event_id": "1",
		"tenant_id": "tenant-1",
		"event_type": "payment_authorized",
		"producer": "manual",
		"event_time": "2026-04-28T00:00:00Z",
		"schema_version": "1",
		"payload": {"amount": 1, "currency": "USD"}
	}`)

	error := processor.ProcessEvent(context.Background(), eventBytes)
	
	if error != nil {
		test.Fatalf("unexpected error: %v", error)
	}

	if len(dataStore.persistedEvents) != 1 {
		test.Fatalf("expected 1 saved event, received: %d", len(dataStore.persistedEvents))
	}
}


func TestProcessEvent_PersistenceError_ReturnsError(test *testing.T) {
	dataStore := &MockDataStore{
		error: errors.New("database failure"),
	}
	processor := newTestProcessor(test, dataStore)

	eventBytes := []byte(`{
		"event_id": "1",
		"tenant_id": "tenant-1",
		"event_type": "TEST",
		"producer": "manual",
		"event_time": "2026-04-28T00:00:00Z",
		"schema_version": "1",
		"payload": {}
	}`)

	error := processor.ProcessEvent(context.Background(), eventBytes)
	
	if error == nil {
		test.Fatal("did not receive expected error")
	}
}

func TestProcessEvent_InvalidEvent_ReturnsError(test *testing.T) {
	dataStore := &MockDataStore{}
	processor := newTestProcessor(test, dataStore)

	eventBytes := []byte(`invalid-json`)

	error := processor.ProcessEvent(context.Background(), eventBytes)
	
	if error == nil {
		test.Fatal("did not receive expected error")
	}

	if len(dataStore.persistedEvents) != 0 {
		test.Fatalf("expected 0 saved events, received %d", len(dataStore.persistedEvents))
	}
}


type MockDataStore struct {
	persistedEvents []event.Event
	error         error
}