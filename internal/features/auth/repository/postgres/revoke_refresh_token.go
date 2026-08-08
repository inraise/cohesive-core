package auth_repository_postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *AuthRepository) RevokeRefreshToken(
	ctx context.Context,
	id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE id = $1;
	`

	if _, err := r.pool.Exec(ctx, query, id, time.Now()); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}
