package auth_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) VerifyEmailByHash(
	ctx context.Context,
	tokenHash string,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH valid_verification AS (
			UPDATE email_verifications
			SET used_at = now()
			WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()
			RETURNING user_id
		)
		UPDATE users
		SET is_verified = true, updated_at = now()
		FROM valid_verification
		WHERE users.id = valid_verification.user_id
		RETURNING
			users.id, users.version, users.email, users.password_hash, users.is_verified,
			users.first_name, users.last_name, users.age, users.created_at, users.updated_at;
	`

	row := r.pool.QueryRow(ctx, query, tokenHash)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.IsVerified,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Age,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf(
				"verification token is invalid, expired, or already used: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		return core_domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Email,
		userModel.PasswordHash,
		userModel.FirstName,
		userModel.LastName,
		userModel.Age,
		userModel.IsVerified,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	), nil
}
