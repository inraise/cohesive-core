package auth_repository_postgres

import (
	"context"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
)

func (r *AuthRepository) CreateEmailVerification(
	ctx context.Context,
	verification core_domain.EmailVerification,
) (core_domain.EmailVerification, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO email_verifications (user_id, email, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, email, token_hash, expires_at, used_at, created_at;
	`

	row := r.pool.QueryRow(ctx, query,
		verification.UserID,
		verification.Email,
		verification.TokenHash,
		verification.ExpiresAt,
		verification.CreatedAt,
	)

	var model EmailVerificationModel
	err := row.Scan(
		&model.ID,
		&model.UserID,
		&model.Email,
		&model.TokenHash,
		&model.ExpiresAt,
		&model.UsedAt,
		&model.CreatedAt,
	)
	if err != nil {
		return core_domain.EmailVerification{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.EmailVerification{
		ID:        model.ID,
		UserID:    model.UserID,
		Email:     model.Email,
		TokenHash: model.TokenHash,
		ExpiresAt: model.ExpiresAt,
		UsedAt:    model.UsedAt,
		CreatedAt: model.CreatedAt,
	}, nil
}
