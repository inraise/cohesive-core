package auth_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *AuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, email, password_hash, first_name, last_name, age, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	row := r.pool.QueryRow(ctx, query, email)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Age,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf("user with email %q: %w", email, core_errors.ErrNotFound)
		}

		return core_domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := core_domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Email,
		userModel.PasswordHash,
		userModel.FirstName,
		userModel.LastName,
		userModel.Age,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)

	return userDomain, nil
}
