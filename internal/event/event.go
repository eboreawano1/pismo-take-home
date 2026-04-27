package event

import (
	"time"
	"encoding/json"
)

type Event struct {
	EventId       string          `json:"event_id"`
	EventType     string          `json:"event_type"`
	TenantId      string          `json:"tenant_id"`
	Producer      string          `josn:"producer`
	EventTime     time.Time       `json:"event_time"`
	Payload       json.RawMessage `json:"payload"`
	SchemaVersion string          `json:"schema_version"`
}