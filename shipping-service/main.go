package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"com.ecommerce/pkg/common/config"
	"com.ecommerce/pkg/common/database"
	"com.ecommerce/pkg/common/middleware"
	"com.ecommerce/shipping-service/internal/handler"
	"com.ecommerce/shipping-service/internal/repository"
	"com.ecommerce/shipping-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func buildPostgresDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	user := config.GetEnv("POSTGRES_USER", "postgres")
	pass := config.GetEnv("POSTGRES_PASSWORD", "postgres")
	host := config.GetEnv("POSTGRES_HOST", "localhost")
	port := config.GetEnv("POSTGRES_PORT", "5432")
	db := config.GetEnv("POSTGRES_DB", "shippingservice")

	if springURL := os.Getenv("SPRING_DATASOURCE_URL"); springURL != "" {
		cleanURL := strings.TrimPrefix(springURL, "jdbc:postgresql://")
		parts := strings.Split(cleanURL, "/")
		if len(parts) >= 2 {
			hostPort := parts[0]
			db = parts[1]
			hpParts := strings.Split(hostPort, ":")
			if len(hpParts) == 2 {
				host = hpParts[0]
				port = hpParts[1]
			} else {
				host = hostPort
			}
		}
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
}

func main() {
	port := config.GetEnv("PORT", "8087")
	dsn := buildPostgresDSN()

	log.Printf("[Shipping Service] Starting service on port %s...", port)
	ctx := context.Background()

	pool, err := database.ConnectPostgres(ctx, dsn)
	if err != nil {
		log.Printf("[Shipping Service] Warning: Failed to connect to Postgres: %v", err)
	} else {
		defer pool.Close()
		log.Println("[Shipping Service] Connected to PostgreSQL successfully")
	}

	repo := repository.NewShippingRepository(pool)
	repo.InitSchema(ctx)

	svc := service.NewShippingService(repo)
	h := handler.NewShippingHandler(svc)

	app := fiber.New(fiber.Config{
		AppName:               "Ecommerce Shipping Service (Golang)",
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

	log.Println("[Shipping Service] Shutting down server gracefully...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	log.Println("[Shipping Service] Server exited successfully.")
}
