package auth_repository_postgres

import (
	"context"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
)

func (r *AuthRepository) CreateRefreshToken(
	ctx context.Context,
	token core_domain.RefreshToken,
) (core_domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at;
	`

	row := r.pool.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
	)

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
