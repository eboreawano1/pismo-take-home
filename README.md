# Pismo Event Processor

This service is an event processor that consumes events from multiple producers, validates the event and payload, persists the event in a postgres database, and persists a delivery target for future delivery services.

## Design summary

This system uses an at-least-once processing model to ensure all events are processed and has idempotency checks on persistence as a safeguard to prevent events from being persisted multiple times. Kafka offsets are committed only after the event is successfully persisted to the database. This allows events to attempt to be processed multiple times in the event of transient errors. Idempotency checks on persistence make these attempts safe.

## Architecture

Producer --> Kafka topic(events) --> Event Processor --> PostgreSQL --> Delivery Service


## Processing behavior

There are 4 cases that the design explicitly handles:
1) valid events persist with the READY_TO_DELIVER status and are routed appropriately
2) Invalid events are still persisted with the INVALID_EVENT status to ensure no data is lost. The validation errors are also persisted for visibility to help with debugging and resolving validation issues.
3) A "valid" but unroutable event was introduced to test the routing behavior explicity. In practice, valid events should have a delivery target. In the case that the delivery target isn't handled by the router the event is persisted with a PROCESSING_ERROR.
4) Duplicate events are handled by using the event_id as the idempotency key. I intentionally chose not to update existing rows on duplicate event_Ids for simplicity. This design assumes that the first event persisted is always the most correct, but this assumption can be revisited depending on the context with the surrounding system.

## Event shape

Events should follow this structure:

```json
{
  "event_id": "event-1",
  "event_type": "payment_authorized",
  "event_time": "2026-04-28T00:00:00Z",
  "tenant_id": "tenant-1",
  "producer": "payments-api",
  "schema_version": "1",
  "payload": {
    "amount": 1000,
    "currency": "USD"
  }
}
```

The event is validated explicitly with business logic and the payload is validated with a stored schema in this format:
schemas/payment_authorized/1.json


## Routing
payment_authorized -> analytics
user_sign_up -> notifications


## Setup requirements

- Docker
- Docker Compose
- Go 1.23+
- jq


## Run the full environment

From the root folder of the project:

```bash
docker compose up --build -d
```

This starts PostgresSQL, Kafka, Zookeeper (not used explicitly just as Kafka setup), and the event processor. The processed_events table is created with a .sql script in the /db folder.

## Create the Kafka topic

```bash
docker compose exec kafka kafka-topics \
  --bootstrap-server kafka:9092 \
  --create \
  --if-not-exists \
  --topic events \
  --partitions 1 \
  --replication-factor 1
```

## Publish a sample event

```bash
jq -c . test/events/payment_authorized.json | docker compose exec -T kafka kafka-console-producer \
  --bootstrap-server kafka:9092 \
  --topic events
```

## Verify persisted events

Open Postgres:

```bash
docker compose exec postgres psql -U username -d pismo-database
```

Then run:

```sql
SELECT
    event_id,
    tenant_id,
    event_type,
    producer,
    event_time,
    schema_version,
    status,
    delivery_target,
    validation_errors,
    create_time
FROM processed_events
ORDER BY create_time DESC;
```

For the sample event, you should see a row with:

```text
status = READY_TO_DELIVER
delivery_target = analytics
```

## Run tests

```bash
go test ./...
```

Test cases:

1) valid events -> persisted as READY_TO_DELIVER
2) invalid events -> persisted as INVALID_EVENT with validation errors
3) unroutable events -> persisted with PROCESSING_ERROR


## Tradeoffs

This design intentionally biases towards correctness and simplicity

The design includes:

- A Kafka-based reactive or pull model for event consumption. This supports back pressure by allowing consumers to work at their own pace as opposed to adapting to the pace of the broker.
- Manual offset commits instead of automatic offset commits. This prevents data loss by enforcing successful persistence is a prerequisite to an offset commit.
- Payload validation using JSON Schema. This allows us to validate events strictly, dynamically, and extensibly. This model makes it easy to adopt new schemas for event types over time and preserves a version history to make rollbacks easier as well. Advantageous for both rapid iteration and backwards compatibility.
- Invalid events persisted. This prevents data loss and provides a durable record of processing failure for future investigation.
- persisted processing errors for unroutable events. Same advantages as discussed above.
- Idempotent persistence using the event_id as the idempotency key. Allows for safe retries on events for transient errors.
- Docker-based services for reproducibility and testing advantages (end-to-end/Integration testing)

Not included:

- Dynamic routing 
- Dead-letter topic / bounded retries / backoff strategy
- Concurrent processing with order correctness enforcement
- Graceful shutdown strategy

These considerations are valid for a production ready service but not necessary for take-home scope. I biased toward keeping the implementation as simple as possible so as to not over-engineer in anticipation of requirements that were not immediately present.


## Potential improvements

- Graceful shutdown to stop polling for new kafka messages and complete processing of all in-flight messages before closing Kafka and db connections.
- Partition-aware and bounded concurrency so message offsets aren't committed out of order per-partitiion. Concurrency improves throughput but would be bounded to provide backpressure to keep consumers from being overwhelmed.
- Explicit retry and dead-letter handling to prevent "poison messages" or infinite retries on non-transient errors. Retrying with a backoff as well to prevent consumers from being overwhelmed.
- Schema and routing management using configuration sources instead of local files for dynamic changes to routing or schemas.
- Integration/end-to-end test scripts to dynamically verify the happy path for automated regression testing.