package auth_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) GetRefreshTokenByHash(
	ctx context.Context,
	tokenHash string,
) (core_domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1;
	`

	row := r.pool.QueryRow(ctx, query, tokenHash)

	var model RefreshTokenModel
	err := row.Scan(
		&model.ID,
		&model.UserID,
		&model.TokenHash,
		&model.ExpiresAt,
		&model.RevokedAt,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.RefreshToken{}, fmt.Errorf("refresh token: %w", core_errors.ErrNotFound)
		}

		return core_domain.RefreshToken{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.RefreshToken{
		ID:        model.ID,
		UserID:    model.UserID,
		TokenHash: model.TokenHash,
		ExpiresAt: model.ExpiresAt,
		RevokedAt: model.RevokedAt,
		CreatedAt: model.CreatedAt,
	}, nil
}
