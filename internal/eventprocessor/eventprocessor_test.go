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

func (dataStore *MockDataStore) Save(
		ctx context.Context, 
		event event.Event, 
		status string, 
		deliveryTarget string,
		validationErrors []string,
	) error {
	if dataStore.error != nil {
		return dataStore.error
	}

	dataStore.persistedEvents = append(dataStore.persistedEvents, event)
	dataStore.statuses = append(dataStore.statuses, status)
	dataStore.validationErrors = append(dataStore.validationErrors, validationErrors)
	dataStore.deliveryTargets = append(dataStore.deliveryTargets, deliveryTarget)

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

	if dataStore.validationErrors[0] != nil {
		test.Fatalf("expected no validation errors but found: %v", dataStore.validationErrors[0])
	}

	if dataStore.deliveryTargets[0] != "analytics" {
		test.Fatalf("expected analytics but found: %s", dataStore.deliveryTargets[0])
	}

	if dataStore.statuses[0] != event.ReadyToDeliverStatus {
		test.Fatalf("expected READY_TO_DELIVER status but found: %s", dataStore.statuses[0])
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
		"event_type": "payment_authorized",
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

func TestProcessEvent_InvalidEvent_InvalidEventPersisted(test *testing.T) {
	dataStore := &MockDataStore{}
	processor := newTestProcessor(test, dataStore)

	eventBytes := []byte(`{
		"event_id": "",
		"tenant_id": "tenant-1",
		"event_type": "payment_authorized",
		"producer": "payments-api",
		"event_time": "2026-04-28T20:00:00Z",
		"schema_version": "1",
		"payload": {}
	}`)

	error := processor.ProcessEvent(context.Background(), eventBytes)

	if error != nil {
		test.Fatalf("unexpected error: %v", error)
	}

	if len(dataStore.persistedEvents) != 1 {
		test.Fatalf("expected 1 saved event, found: %d", len(dataStore.persistedEvents))
	}

	if dataStore.statuses[0] != event.InvalidEventStatus {
		test.Fatalf("expected INVALID status, found %s", dataStore.statuses[0])
	}

	if dataStore.validationErrors[0] == nil {
		test.Fatal("missing expected validation errors")
	}
}

func TestProcess_UnroutableEvent_PersistedWithProcessingError(test *testing.T) {
	dataStore := &MockDataStore{}
	processor := newTestProcessor(test, dataStore)

	eventBytes := []byte(`{
		"event_id": "event-1",
		"tenant_id": "tenant-1",
		"event_type": "unrouted_event",
		"producer": "new-api",
		"event_time": "2026-04-28T20:00:00Z",
		"schema_version": "1",
		"payload": {}
	}`)

	error := processor.ProcessEvent(context.Background(), eventBytes)

	if error != nil {
		test.Fatalf("received unexpected error:  %v", error)
	}

	if dataStore.statuses[0] != event.ProcessingErrorStatus {
		test.Fatalf("expected PROCESSING_ERROR status but found: %s", dataStore.statuses[0])
	}

	if dataStore.deliveryTargets[0] != "" {
		test.Fatalf("expected missing delivery target but found:  %s", dataStore.deliveryTargets[0])
	}

	if len(dataStore.persistedEvents) != 1 {
		test.Fatalf("expected 1 saved event but found: %d", len(dataStore.persistedEvents))
	}
}

type MockDataStore struct {
	persistedEvents []event.Event
	error         error
	statuses    []string
	validationErrors     [][]string
	deliveryTargets []string
}