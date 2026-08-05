package auth_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	"context"
)

type AuthService struct {
	authRepository AuthRepository
	tokenManager   *core_jwt.TokenManager
}

type AuthRepository interface {
	CreateUser(
		ctx context.Context,
		user core_domain.User,
	) (core_domain.User, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (core_domain.User, error)
}

func NewAuthService(
	authRepository AuthRepository,
	tokenManager *core_jwt.TokenManager,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		tokenManager:   tokenManager,
	}
}
