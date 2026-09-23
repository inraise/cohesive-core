package users_repository_postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"

	"github.com/google/uuid"
)

func (r *UsersRepository) PatchMe(
	ctx context.Context,
	id uuid.UUID,
	user core_domain.User,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE users
	SET 
		email=$1,
		password_hash=$2,
		is_verified=$3,
		first_name=$4,
		last_name=$5,
		age=$6,
		updated_at=$7,
		version=version+1
	WHERE id=$8 AND version=$9
	RETURNING
		id,
		version,
		email,
		password_hash,
		is_verified,
		first_name,
		last_name,
		age,
		created_at,
		updated_at;`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.IsVerified,
		user.FirstName,
		user.LastName,
		user.Age,
		time.Now(),
		user.ID,
		user.Version,
	)

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
				"user with id=%q concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
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
		userModel.IsVerified,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)

	return userDomain, nil
}
