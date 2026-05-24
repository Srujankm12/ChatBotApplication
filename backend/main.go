package main

import (
	"context"
	"log"
	"os"
	"time"

	"chatbot/internal/handler"
	"chatbot/internal/repository"
	"chatbot/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load("../.env")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}
	log.Println("connected to MongoDB")

	db := client.Database("chatbot")

	chatRepo := repository.NewChatRepository(db)
	logRepo := repository.NewLogRepository(db)

	chatSvc := service.NewChatService(chatRepo)
	metricsSvc := service.NewMetricsService(logRepo)

	chatHandler := handler.NewChatHandler(chatSvc)
	logHandler := handler.NewLogHandler(logRepo)
	metricsHandler := handler.NewMetricsHandler(metricsSvc)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	e.POST("/chat", chatHandler.HandleChat)
	e.GET("/chat/:session_id", chatHandler.HandleGetChat)
	e.GET("/sessions", chatHandler.HandleGetSessions)
	e.POST("/logs", logHandler.HandleLogs)
	e.GET("/metrics", metricsHandler.HandleMetrics)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server starting on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
