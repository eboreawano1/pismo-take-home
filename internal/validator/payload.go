package validator


import (
	"fmt"
	"bytes"
	"encoding/json"
	"pismo-take-home/internal/event"
	"pismo-take-home/schemas"
)


func NewPayloadValidator(schemaLoader *schemas.SchemaLoader)  *PayloadValidator {
	return &PayloadValidator{ schemaLoader: schemaLoader }
}

func (validator *PayloadValidator) ValidatePayload(event event.Event) []string {
	schemaPath := fmt.Sprintf("schemas/%s/%s.json", event.EventType, event.SchemaVersion)

	schema, error := validator.schemaLoader.Compile(schemaPath)
	var payload any

	if error != nil {
		return []string{ fmt.Sprintf("Error: schema is invalid: %v", error) }
	}
	
	decodeError := json.NewDecoder(bytes.NewReader(event.Payload)).Decode(&payload)

	if  decodeError != nil {
		return []string{ fmt.Sprintf("Error: payload json is not valid: %v", decodeError) }
	}

	schemaValidationError := schema.Validate(payload)

	if  schemaValidationError != nil {
		return []string{ fmt.Sprintf("Error: schema validation failed for payload: %v", schemaValidationError) }
	}

	return nil
}

type PayloadValidator struct {
	schemaLoader *schemas.SchemaLoader
}