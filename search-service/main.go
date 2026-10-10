package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"com.ecommerce/pkg/common/config"
	commonmw "com.ecommerce/pkg/common/middleware"
	"com.ecommerce/search-service/internal/consumer"
	"com.ecommerce/search-service/internal/handler"
	"com.ecommerce/search-service/internal/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Correlation-Id"},
	}))
	r.Use(commonmw.CorrelationID)
	r.Use(commonmw.UserClaims)

	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Search Service] Shutting down server gracefully...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	log.Println("[Search Service] Server exited successfully.")
}
