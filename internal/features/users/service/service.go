package users_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"

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
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
