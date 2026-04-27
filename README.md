# Pismo Event Processor

Event processor for the Pismo take-home test.

## Setup Requirements

- Docker
- Go 1.23+

## Setup commands

 - setup postgres db
 - create a docker instance for the event processor
 - run the processor to persist a test event
 
```bash
docker compose up -d postgres

docker build -t pismo-eventprocessor .

docker run --rm \
  --network pismo-take-home_default \
  -e "DATABASE_URL=postgres://username:password@postgres:5432/pismo-database?sslmode=disable" \
  pismo-eventprocessor
```