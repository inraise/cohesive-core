package testhelpers

import (
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	core_pool_pgx "cohesive-core/internal/core/repository/postgres/pool/pgx"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	tc_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testUser     = "cohesive_test"
	testPassword = "cohesive_test"
	testDatabase = "cohesive_test"
)

func NewPostgres(t *testing.T) core_pool.Pool {
	t.Helper()

	ctx := context.Background()

	container, err := tc_postgres.Run(ctx,
		"postgres:17-alpine",
		tc_postgres.WithDatabase(testDatabase),
		tc_postgres.WithUsername(testUser),
		tc_postgres.WithPassword(testPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("terminate postgres container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("get container host: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("get container port: %v", err)
	}

	if err := runMigrations(host, mappedPort.Port()); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	pool, err := core_pool_pgx.NewPool(ctx, core_pool_pgx.Config{
		Host:     host,
		Port:     mappedPort.Port(),
		User:     testUser,
		Password: testPassword,
		Database: testDatabase,
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("connect to test postgres: %v", err)
	}

	t.Cleanup(pool.Close)

	return pool
}

func migrationsDir() (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("determine caller info for migrations path")
	}

	return filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations"), nil
}

func runMigrations(host, port string) error {
	dir, err := migrationsDir()
	if err != nil {
		return fmt.Errorf("resolve migrations dir: %w", err)
	}

	dsn := fmt.Sprintf(
		"pgx5://%s:%s@%s:%s/%s?sslmode=disable",
		testUser, testPassword, host, port, testDatabase,
	)

	m, err := migrate.New("file://"+dir, dsn)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
