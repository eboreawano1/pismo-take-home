FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o eventprocessor ./cmd/eventprocessor

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/eventprocessor .
COPY --from=builder /app/schemas ./schemas

CMD ["./eventprocessor"]