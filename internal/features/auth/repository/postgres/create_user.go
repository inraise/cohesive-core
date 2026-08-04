package auth_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"
)

func (r *AuthRepository) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO users (version, email, password_hash, first_name, last_name, age, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, version, email, password_hash, first_name, last_name, age, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		user.Version,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Age,
		user.CreatedAt,
		user.UpdatedAt,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
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
		user.CreatedAt,
		user.UpdatedAt,
	)

	return userDomain, nil
}
