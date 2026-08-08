package auth_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepository AuthRepository
	tokenManager   *core_jwt.TokenManager
	refreshTTL     time.Duration
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

	CreateRefreshToken(
		ctx context.Context,
		token core_domain.RefreshToken,
	) (core_domain.RefreshToken, error)

	GetRefreshTokenByHash(
		ctx context.Context,
		tokenHash string,
	) (core_domain.RefreshToken, error)

	RevokeRefreshToken(
		ctx context.Context,
		id uuid.UUID,
	) error
}

func NewAuthService(
	authRepository AuthRepository,
	tokenManager *core_jwt.TokenManager,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		tokenManager:   tokenManager,
		refreshTTL:     refreshTTL,
	}
}
