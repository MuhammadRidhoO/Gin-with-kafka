package config

import (
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func InitKafkaWriter(topic string) *kafka.Writer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		log.Fatal("KAFKA_BROKER env is required")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}

	log.Println("✅ Kafka writer ready for topic:", topic)
	return writer
}
