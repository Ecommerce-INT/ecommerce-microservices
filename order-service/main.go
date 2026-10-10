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

	"com.ecommerce/order-service/internal/handler"
	"com.ecommerce/order-service/internal/repository"
	"com.ecommerce/order-service/internal/service"
	"com.ecommerce/pkg/common/config"
	"com.ecommerce/pkg/common/database"
	"com.ecommerce/pkg/common/middleware"
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
	db := config.GetEnv("POSTGRES_DB", "orderservice")

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
	port := config.GetEnv("PORT", "8084")
	dsn := buildPostgresDSN()

	log.Printf("[Order Service] Starting service on port %s...", port)
	ctx := context.Background()

	pool, err := database.ConnectPostgres(ctx, dsn)
	if err != nil {
		log.Printf("[Order Service] Warning: Failed to connect to Postgres: %v", err)
	} else {
		defer pool.Close()
		log.Println("[Order Service] Connected to PostgreSQL successfully")
	}

	repo := repository.NewOrderRepository(pool)
	_ = repo.InitSchema(ctx)

	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	app := fiber.New(fiber.Config{
		AppName:               "Ecommerce Order Service (Golang)",
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

	log.Println("[Order Service] Shutting down server gracefully...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	log.Println("[Order Service] Server exited successfully.")
}
