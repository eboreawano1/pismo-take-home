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

jq -c . test/events/payment_authorized.json | docker compose exec -T kafka kafka-console-producer \
  --bootstrap-server kafka:9092 \
  --topic events

docker compose exec postgres psql -U username -d pismo-database

SELECT 
    event_id, 
    tenant_id, 
    event_type, 
    producer, 
    event_time, 
    schema_version, 
    payload, 
    status, 
    create_time
FROM processed_events;
```