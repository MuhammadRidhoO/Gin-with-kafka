package routes

import (
	"go-kafka/internal/controllers"
	"go-kafka/internal/kafka"
	"go-kafka/internal/repositories"
	"go-kafka/internal/service"

	"github.com/gin-gonic/gin"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func SetupRouter(kafkaWriter *kafkago.Writer, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// layering clean architecture
	userRepo := repositories.NewUserRepository(db)
	userUC := service.NewUserUsecase(userRepo, kafka.NewProducer(kafkaWriter))
	userHandler := controllers.NewUserHandler(userUC)

	// routes
	r.POST("/users", userHandler.CreateUser)

	return r
}
