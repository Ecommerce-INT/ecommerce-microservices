package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"com.ecommerce/pkg/common/config"
	"com.ecommerce/pkg/common/middleware"
	"com.ecommerce/search-service/internal/consumer"
	"com.ecommerce/search-service/internal/handler"
	"com.ecommerce/search-service/internal/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	port := config.GetEnv("PORT", "8094")
	esURL := config.GetEnv("ELASTICSEARCH_URL", "http://localhost:9200")
	esUser := config.GetEnv("ELASTICSEARCH_USERNAME", "")
	esPass := config.GetEnv("ELASTICSEARCH_PASSWORD", "")
	kafkaServers := config.GetEnv("KAFKA_SERVERS", "localhost:9092")
	topicName := config.GetEnv("PRODUCT_TOPIC_NAME", "dbproduct.public.product")

	log.Printf("[Search Service] Starting service on port %s...", port)

	esCfg := elasticsearch.Config{
		Addresses: []string{esURL},
		Username:  esUser,
		Password:  esPass,
	}

	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Printf("[Search Service] Warning: Failed to initialize Elasticsearch client: %v", err)
	} else {
		log.Println("[Search Service] Elasticsearch client initialized successfully.")
	}

	svc := service.NewSearchService(esClient)
	h := handler.NewSearchHandler(svc)

	// Start Kafka sync consumer in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	brokers := strings.Split(kafkaServers, ",")
	syncConsumer := consumer.NewProductSyncConsumer(svc, brokers, topicName, "search-group")
	syncConsumer.Start(ctx)

	app := fiber.New(fiber.Config{
		AppName:               "Ecommerce Search Service (Golang)",
		DisableStartupMessage: false,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Correlation-Id",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(middleware.CorrelationID())
	app.Use(middleware.UserClaims())

	h.RegisterRoutes(app)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Search Service] Shutting down server gracefully...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	log.Println("[Search Service] Server exited successfully.")
}
