package users_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"

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
		first_name=$2,
		last_name=$3,
		age=$4,
		updated_at=$5,
		version=version+1
	WHERE id=$6 AND version=$7
	RETURNING
		id,
		version,
		email,
		password_hash,
		first_name,
		last_name,
		age,
		created_at,
		updated_at;`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
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
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Age,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed: %w",
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
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)

	return userDomain, nil
}
