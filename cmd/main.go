package main

import (
	"fmt"
	"go-kafka/config"
	"go-kafka/internal/kafka"
	"go-kafka/internal/routes"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	// DB
	db := config.InitDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get sql.DB from gorm.DB: %v", err)
	}
	defer sqlDB.Close()

	// Kafka
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "user-topic"
	}

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}
	fmt.Println(broker, "Check Broker in main.go")
	fmt.Println(topic, "Check Topic in main.go")

	kafkaWriter := config.InitKafkaWriter(topic)
	defer kafkaWriter.Close()

	// Start Kafka consumer in background
	go kafka.StartConsumer(broker, topic, "user-group")

	// Setup router
	app := routes.SetupRouter(kafkaWriter, db)

	log.Println("🚀 Server running on :8080")
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}

	// keep main alive
	for {
		time.Sleep(time.Second)
	}
}
