package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"apschool/internal/auth"
	"apschool/internal/challenges"
	"apschool/internal/submissions"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	database   = os.Getenv("APSCHOOL_DB_DATABASE")
	password   = os.Getenv("APSCHOOL_DB_PASSWORD")
	username   = os.Getenv("APSCHOOL_DB_USERNAME")
	port       = os.Getenv("APSCHOOL_DB_PORT")
	host       = os.Getenv("APSCHOOL_DB_HOST")
	schema     = os.Getenv("APSCHOOL_DB_SCHEMA")
	serverPort = os.Getenv("PORT")
	sslmode    = os.Getenv("APSCHOOL_DB_SSLMODE")
)

type application struct {
	db          *sql.DB
	logger      *slog.Logger
	auth        *auth.Handler
	challenges  *challenges.Handler
	submissions *submissions.Handler
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	// Run database migrations
	if err := runMigrations(db, "migrations"); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}
	log.Println("Migrations completed successfully")

	app := &application{
		db:          db,
		logger:      logger,
		auth:        auth.NewHandler(auth.NewService(auth.NewRepository(db)), logger),
		challenges:  challenges.NewHandler(challenges.NewService(challenges.NewRepository(db)), logger),
		submissions: submissions.NewHandler(submissions.NewService(submissions.NewRepository(db)), logger),
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", serverPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan bool, 1)
	go gracefulShutdown(server, done)

	log.Printf("Starting server on port %s", server.Addr)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func openDB() (*sql.DB, error) {

	if sslmode == "" {
		sslmode = "disable"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s", username, password, host, port, database, sslmode, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil

}

// runMigrations executes all SQL migration files in order.
// It reads files from the migrations directory and executes only the "Up" portion.
func runMigrations(db *sql.DB, migrationsPath string) error {
	// Create migrations tracking table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Read migration files
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute each migration
	for _, fileName := range sqlFiles {
		// Check if migration was already applied
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", fileName).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", fileName, err)
		}
		if exists {
			log.Printf("Migration %s already applied, skipping", fileName)
			continue
		}

		// Read migration file
		content, err := os.ReadFile(filepath.Join(migrationsPath, fileName))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		// Extract only the "Up" portion (before -- +goose Down)
		upSQL := extractUpMigration(string(content))
		if upSQL == "" {
			log.Printf("Warning: no Up migration found in %s, skipping", fileName)
			continue
		}

		// Execute migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", fileName, err)
		}

		if _, err := tx.Exec(upSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		// Record migration as applied
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", fileName); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", fileName, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", fileName, err)
		}

		log.Printf("Applied migration: %s", fileName)
	}

	return nil
}

// extractUpMigration extracts the SQL between "-- +goose Up" and "-- +goose Down"
func extractUpMigration(content string) string {
	lines := strings.Split(content, "\n")
	var upLines []string
	inUp := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "-- +goose Up") {
			inUp = true
			continue
		}
		if strings.Contains(trimmed, "-- +goose Down") {
			break
		}
		if inUp {
			upLines = append(upLines, line)
		}
	}

	return strings.TrimSpace(strings.Join(upLines, "\n"))
}
