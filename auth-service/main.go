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

	"com.ecommerce/auth-service/internal/handler"
	"com.ecommerce/auth-service/internal/repository"
	"com.ecommerce/auth-service/internal/service"
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
	db := config.GetEnv("POSTGRES_DB", "authservice")

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
	port := config.GetEnv("PORT", "8088")
	dsn := buildPostgresDSN()

	log.Printf("[Auth Service] Starting service on port %s...", port)
	ctx := context.Background()

	pool, err := database.ConnectPostgres(ctx, dsn)
	if err != nil {
		log.Printf("[Auth Service] Warning: Failed to connect to Postgres: %v", err)
	} else {
		defer pool.Close()
		log.Println("[Auth Service] Connected to PostgreSQL successfully")
	}

	repo := repository.NewUserRepository(pool)
	_ = repo.InitSchema(ctx)

	kcCfg := service.KeycloakConfig{
		ServerURL:               config.GetEnv("KEYCLOAK_SERVER_URL", "http://localhost:8080"),
		PublicServerURL:         config.GetEnv("KEYCLOAK_PUBLIC_SERVER_URL", "http://keycloak.ecommerce.local"),
		Realm:                   config.GetEnv("KEYCLOAK_REALM", "ecommerce"),
		ClientID:                config.GetEnv("KEYCLOAK_CLIENT_ID", "ecommerce-client"),
		ClientSecret:            config.GetEnv("KEYCLOAK_CLIENT_SECRET", ""),
		AdminRealm:              config.GetEnv("KEYCLOAK_ADMIN_REALM", "master"),
		AdminClientID:           config.GetEnv("KEYCLOAK_ADMIN_CLIENT_ID", "admin-cli"),
		AdminUsername:           config.GetEnv("KEYCLOAK_ADMIN_USERNAME", "admin"),
		AdminPassword:           config.GetEnv("KEYCLOAK_ADMIN_PASSWORD", "admin"),
		BackendCallbackURL:      config.GetEnv("SSO_BACKEND_CALLBACK_URL", "http://api.ecommerce.local/api/v1/auth/callback"),
		DefaultFrontendRedirect: config.GetEnv("SSO_DEFAULT_FRONTEND_REDIRECT", "http://ecommerce.local/auth/callback"),
		Scope:                   config.GetEnv("SSO_SCOPE", "openid profile email"),
	}

	kcClient := service.NewKeycloakClient(kcCfg)
	ssoStore := service.NewSsoSessionStore()
	userSvc := service.NewUserService(repo, kcClient)
	h := handler.NewAuthHandler(userSvc, kcClient, ssoStore, kcCfg.DefaultFrontendRedirect)

	app := fiber.New(fiber.Config{
		AppName:               "Ecommerce Auth Service (Golang)",
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

	log.Println("[Auth Service] Shutting down server gracefully...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	log.Println("[Auth Service] Server exited successfully.")
}
