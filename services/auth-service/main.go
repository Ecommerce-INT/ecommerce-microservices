package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"com.ecommerce/auth-service/internal/handler"
	"com.ecommerce/auth-service/internal/repository"
	"com.ecommerce/auth-service/internal/service"
	"com.ecommerce/shared/config"
	"com.ecommerce/shared/database"
	commonmw "com.ecommerce/shared/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	log.Println("[Auth Service] Shutting down server gracefully...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	log.Println("[Auth Service] Server exited successfully.")
}
