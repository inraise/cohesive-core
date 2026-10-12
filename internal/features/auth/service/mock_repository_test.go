package auth_service_test

import (
	"context"

	core_domain "cohesive-core/internal/core/domain"

	"github.com/google/uuid"
)

type mockAuthRepository struct {
	createUserFunc func(ctx context.Context, user core_domain.User) (core_domain.User, error)

	getUserByEmailFunc func(ctx context.Context, email string) (core_domain.User, error)

	getUserByIDFunc func(ctx context.Context, id uuid.UUID) (core_domain.User, error)

	createRefreshTokenFunc func(
		ctx context.Context,
		token core_domain.RefreshToken,
	) (core_domain.RefreshToken, error)

	getRefreshTokenByHashFunc func(ctx context.Context, tokenHash string) (core_domain.RefreshToken, error)

	revokeRefreshTokenFunc func(ctx context.Context, id uuid.UUID) error

	createEmailVerificationFunc func(
		ctx context.Context,
		verification core_domain.EmailVerification,
	) (core_domain.EmailVerification, error)

	verifyEmailByHashFunc func(ctx context.Context, tokenHash string) (core_domain.User, error)
}

func (m *mockAuthRepository) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	if m.createUserFunc == nil {
		panic("mockAuthRepository.CreateUser called but createUserFunc is not set")
	}

	return m.createUserFunc(ctx, user)
}

func (m *mockAuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (core_domain.User, error) {
	if m.getUserByEmailFunc == nil {
		panic("mockAuthRepository.GetUserByEmail called but getUserByEmailFunc is not set")
	}

	return m.getUserByEmailFunc(ctx, email)
}

func (m *mockAuthRepository) CreateRefreshToken(
	ctx context.Context,
	token core_domain.RefreshToken,
) (core_domain.RefreshToken, error) {
	if m.createRefreshTokenFunc == nil {
		panic("mockAuthRepository.CreateRefreshToken called but createRefreshTokenFunc is not set")
	}

	return m.createRefreshTokenFunc(ctx, token)
}

func (m *mockAuthRepository) GetRefreshTokenByHash(
	ctx context.Context,
	tokenHash string,
) (core_domain.RefreshToken, error) {
	if m.getRefreshTokenByHashFunc == nil {
		panic("mockAuthRepository.GetRefreshTokenByHash called but getRefreshTokenByHashFunc is not set")
	}

	return m.getRefreshTokenByHashFunc(ctx, tokenHash)
}

func (m *mockAuthRepository) RevokeRefreshToken(
	ctx context.Context,
	id uuid.UUID,
) error {
	if m.revokeRefreshTokenFunc == nil {
		panic("mockAuthRepository.RevokeRefreshToken called but revokeRefreshTokenFunc is not set")
	}

	return m.revokeRefreshTokenFunc(ctx, id)
}

func (m *mockAuthRepository) CreateEmailVerification(
	ctx context.Context,
	verification core_domain.EmailVerification,
) (core_domain.EmailVerification, error) {
	if m.createEmailVerificationFunc == nil {
		panic("mockAuthRepository.CreateEmailVerification called but createEmailVerificationFunc is not set")
	}

	return m.createEmailVerificationFunc(ctx, verification)
}

func (m *mockAuthRepository) VerifyEmailByHash(
	ctx context.Context,
	tokenHash string,
) (core_domain.User, error) {
	if m.verifyEmailByHashFunc == nil {
		panic("mockAuthRepository.VerifyEmailByHash called but verifyEmailByHashFunc is not set")
	}

	return m.verifyEmailByHashFunc(ctx, tokenHash)
}

func (m *mockAuthRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (core_domain.User, error) {
	if m.getUserByIDFunc == nil {
		panic("mockAuthRepository.GetUserByID called but getUserByIDFunc is not set")
	}

	return m.getUserByIDFunc(ctx, id)
}
