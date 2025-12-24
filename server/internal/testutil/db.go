package testutil

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
	"time"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	DB *sql.DB
	Container testcontainers.Container
}

func SetupTestDB(ctx context.Context) (*TestDB, error) {
	container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		container.Terminate(ctx)
		return nil, err
	}

	if err := applyMigrations(db); err != nil {
		db.Close()
		container.Terminate(ctx)
		return nil, err
	}

	return &TestDB{
		DB: db,
		Container: container,
	}, nil
}

func applyMigrations(db *sql.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "migrations")
	goose.SetDialect("postgres")
	return goose.Up(db, migrationsDir)
}

func (tdb *TestDB) Teardown(ctx context.Context) error {
	if tdb.DB != nil {
		tdb.DB.Close()
	}

	if tdb.Container != nil {
		return tdb.Container.Terminate(ctx)
	}

	return nil
}

func (tdb *TestDB) TruncateTables(t *testing.T) {
	t.Helper()
	_, err := tdb.DB.Exec("TRUNCATE users CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}
