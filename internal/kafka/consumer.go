package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(broker, topic, groupID string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
	defer reader.Close()

	log.Println("🟢 Kafka consumer started for topic:", topic)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}
		log.Printf("📩 Received message: %s", string(m.Value))
	}
}
