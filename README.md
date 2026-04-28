# Pismo Event Processor

Event processor for the Pismo take-home test.

## Setup Requirements

- Docker
- Go 1.23+

## Setup commands

 - setup postgres db, and processor/consumer containers
 - create an event kafka topic
 - publish an event
 - query the db to verify that it persisted

```bash
docker compose up --build

docker compose exec kafka kafka-topics \
 --bootstrap-server kafka:9092 \
 --create \
 --if-not-exists \
 --topic events \
 --partitions 1 \
 --replication-factor 1

docker compose exec -T kafka kafka-console-producer \
 --bootstrap-server kafka:9092 \
 --topic events

{
    "event_id":"event-1",
    "tenant_id":"tenant-1",
    "event_type":"payment_authorized",
    "producer":"payments-api",
    "event_time":"2026-04-27T20:00:00Z",
    "schema_version":"1",
    "payload":{"amount":1000,"currency":"USD"}
}

docker compose exec postgres psql -U username -d pismo-database

SELECT event_id, tenant_id, event_type, status, create_time
FROM processed_events;
```