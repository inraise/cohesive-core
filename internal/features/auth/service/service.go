package auth_service

import (
	"context"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	core_mailer "cohesive-core/internal/core/mailer"
	core_verification "cohesive-core/internal/core/verification"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepository AuthRepository
	tokenManager   *core_jwt.TokenManager
	refreshTTL     time.Duration

	mailer       core_mailer.Mailer
	verification core_verification.Config
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

	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
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

	CreateEmailVerification(
		ctx context.Context,
		verification core_domain.EmailVerification,
	) (core_domain.EmailVerification, error)

	VerifyEmailByHash(
		ctx context.Context,
		tokenHash string,
	) (core_domain.User, error)
}

func NewAuthService(
	authRepository AuthRepository,
	tokenManager *core_jwt.TokenManager,
	refreshTTL time.Duration,
	mailer core_mailer.Mailer,
	verification core_verification.Config,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		tokenManager:   tokenManager,
		refreshTTL:     refreshTTL,
		mailer:         mailer,
		verification:   verification,
	}
}
