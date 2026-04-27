CREATE TABLE processed_events (
    event_id VARCHAR(255) PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    producer TEXT NOT NULL,
    status VARCHAR(255) NOT NULL,
    schema_version VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    delivery_target TEXT,
    validation_errors JSONB,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    event_time TIMESTAMPTZ NOT NULL
)