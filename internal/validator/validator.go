package validator

import (
	"fmt"
	"strings"
	"pismo-take-home/internal/event"
)

func Validate(event event.Event)  []string {
	var errors []string

	if strings.TrimSpace(event.EventId) == "" {
		errors = append(errors, "Error: event_id is missing")
	}

	if strings.TrimSpace(event.EventType) == "" {
		errors = append(errors, "Error: event_type is missing")
	}

	if strings.TrimSpace(event.TenantId) == "" {
		errors = append(errors, "Error: tenant_id is missing")
	}

	if strings.TrimSpace(event.Producer) == "" {
		errors = append(errors, "Error: producer is missing")
	}

	if event.EventTime.IsZero() {
		errors = append(errors, "Error: event time is invalid")
	}

	if strings.TrimSpace(event.SchemaVersion) == "" {
		errors = append(errors, "Error: schema_version is missing")
	}

	return errors
}