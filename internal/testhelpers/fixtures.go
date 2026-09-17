package testhelpers

import (
	"context"
	"testing"
	"time"

	core_pool "cohesive-core/internal/core/repository/postgres/pool"

	"github.com/google/uuid"
)

func InsertUser(t *testing.T, pool core_pool.Pool, email string) uuid.UUID {
	t.Helper()

	id := uuid.New()
	now := time.Now()

	_, err := pool.Exec(context.Background(), `
		INSERT INTO users (id, version, email, password_hash, first_name, created_at, updated_at)
		VALUES ($1, 1, $2, 'test-hash-not-a-real-bcrypt-hash', 'Test', $3, $3);
	`, id, email, now)
	if err != nil {
		t.Fatalf("insert test user: %v", err)
	}

	return id
}
