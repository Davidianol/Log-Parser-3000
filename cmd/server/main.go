package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"log_parser3000/internal/handler"
	"log_parser3000/internal/parser"
	"log_parser3000/internal/repository/postgres"
	"log_parser3000/internal/service"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	initLogger()

	databaseURL := getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/logparser?sslmode=disable")
	port := getenv("PORT", "8080")
	dataDir := getenv("DATA_DIR", "./data")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		slog.Error("open db failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("ping db failed", "err", err)
		os.Exit(1)
	}

	if err := runMigrations(db); err != nil {
		slog.Error("migrations failed", "err", err)
		os.Exit(1)
	}

	repo := postgres.New(db)
	p := parser.NewMainParser(repo, dataDir)
	svc := service.New(repo, p)
	h := handler.New(svc)
	router := handler.NewRouter(h)

	slog.Info("server started", "port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

func runMigrations(db *sql.DB) error {
	driver, err := pgmigrate.WithInstance(db, &pgmigrate.Config{})
	if err != nil {
		return fmt.Errorf("migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}
	slog.Info("migrations ok")
	return nil
}

func initLogger() {
	level := slog.LevelInfo
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
