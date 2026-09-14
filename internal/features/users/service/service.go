package users_service

import (
	"context"

	core_domain "cohesive-core/internal/core/domain"

	"github.com/google/uuid"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (core_domain.User, error)

	PatchMe(
		ctx context.Context,
		id uuid.UUID,
		user core_domain.User,
	) (core_domain.User, error)

	DeleteMe(
		ctx context.Context,
		id uuid.UUID,
	) error
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
