package consumer


import (
	"fmt"
	"context"
	"github.com/segmentio/kafka-go"
)

func NewKafkaConsumer(topic string, groupId string, brokers []string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{		
		Topic:       topic,
		GroupID:     groupId,
		Brokers:     brokers,
		StartOffset: kafka.FirstOffset,
	})

	return &KafkaConsumer{ reader: reader }
}

func (consumer *KafkaConsumer) ConsumeMessage(ctx context.Context) (*Message, error) {
	message, error := consumer.reader.FetchMessage(ctx)

	if error != nil {
		return nil, fmt.Errorf("Error fetching message: %w", error)
	}

	return &Message{
		ByteValue: message.Value,
		raw:   message,
	}, nil
}

func (consumer *KafkaConsumer) CommitMessage(ctx context.Context, message *Message) error {
	commitError := consumer.reader.CommitMessages(ctx, message.raw)

	if commitError != nil {
		return fmt.Errorf("Error committing message: %w", commitError)
	}

	return nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}


type Message struct {
	raw   kafka.Message
	ByteValue []byte
}

type KafkaConsumer struct {
	reader *kafka.Reader
}