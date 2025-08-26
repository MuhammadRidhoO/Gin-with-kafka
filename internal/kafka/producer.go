package kafka

import (
	"context"
	"encoding/json"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, key string, payload any) error
}

type producer struct {
	writer *kafkago.Writer
}

func NewProducer(writer *kafkago.Writer) Producer {
	return &producer{writer: writer}
}

func (p *producer) Publish(ctx context.Context, key string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(key),
		Value: b,
		Headers: []kafkago.Header{
			{Key: "content-type", Value: []byte("application/json")},
		},
	})
	if err != nil {
		log.Println("Failed to publish:", err)
		return err
	}
	log.Println("✅ Message published:", string(b))
	return nil
}
