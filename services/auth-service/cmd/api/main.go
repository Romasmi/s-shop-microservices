package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getDbUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	dbUrl := getDbUrl()
	slog.Info("Connecting to database", "url", dbUrl)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		slog.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Run migrations
	basePath := os.Getenv("APP_BASE_PATH")
	if basePath == "" {
		basePath = "."
	}
	migrationsPath := filepath.Join(basePath, "migrations")
	absPath, _ := filepath.Abs(migrationsPath)
	slog.Info("Running migrations", "path", absPath)
	m, err := migrate.New(
		"file://"+absPath,
		dbUrl,
	)
	if err != nil {
		slog.Error("unable to create migrations driver", "error", err)
	} else {
		err = m.Up()
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				slog.Info("no migrations to apply")
			} else {
				slog.Error("error while running up migrations", "error", err)
				os.Exit(1)
			}
		} else {
			slog.Info("migrations applied successfully")
		}
		m.Close()
	}

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"OK"}`)
	}).Methods(http.MethodGet)

	router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		_, err := dbPool.Exec(r.Context(), "INSERT INTO auth_logs (user_id, login) VALUES ($1, $2)", "1", "admin")
		if err != nil {
			slog.Error("failed to log auth to db", "error", err)
		}

		w.Header().Set("X-Auth-User-ID", "1")
		w.Header().Set("X-Auth-User-Login", "admin")
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		slog.Info("Auth service starting", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	slog.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error during shutdown", "error", err)
	}

	slog.Info("Server stopped")
}
