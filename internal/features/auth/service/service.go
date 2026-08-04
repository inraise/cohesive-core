package auth_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
)

type AuthService struct {
	authRepository AuthRepository
}

type AuthRepository interface {
	CreateUser(
		ctx context.Context,
		user core_domain.User,
	) (core_domain.User, error)
}

func NewAuthService(authRepository AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}
