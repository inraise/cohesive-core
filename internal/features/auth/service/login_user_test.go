package auth_service_test

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func mustHashPassword(t *testing.T, password string) string {
	t.Helper()

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	return string(hashed)
}

func TestAuthService_LoginUser(t *testing.T) {
	const correctPassword = "correct-horse-battery-staple"

	existingUser := core_domain.NewUser(
		uuid.New(),
		1,
		"user@example.com",
		mustHashPassword(t, correctPassword),
		"Nikita",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)

	tests := []struct {
		name      string
		email     string
		password  string
		repo      *mockAuthRepository
		wantErr   bool
		wantErrIs error
	}{
		{
			name:     "успешный логин выдаёт access и refresh токены",
			email:    existingUser.Email,
			password: correctPassword,
			repo: &mockAuthRepository{
				getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
					return existingUser, nil
				},
				createRefreshTokenFunc: func(
					_ context.Context,
					token core_domain.RefreshToken,
				) (core_domain.RefreshToken, error) {
					token.ID = uuid.New()

					return token, nil
				},
			},
			wantErr: false,
		},
		{
			name:     "неверный пароль маскируется под ErrInvalidArgument",
			email:    existingUser.Email,
			password: "totally-wrong-password",
			repo: &mockAuthRepository{
				getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
					return existingUser, nil
				},
			},
			wantErr:   true,
			wantErrIs: core_errors.ErrInvalidArgument,
		},
		{
			name:     "несуществующий email маскируется под ту же ошибку, что и неверный пароль",
			email:    "ghost@example.com",
			password: correctPassword,
			repo: &mockAuthRepository{
				getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
					return core_domain.User{}, fmt.Errorf("user lookup: %w", core_errors.ErrNotFound)
				},
			},
			wantErr:   true,
			wantErrIs: core_errors.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestAuthService(t, tt.repo)

			result, err := service.LoginUser(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)

				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.NotEmpty(t, result.AccessToken)
			assert.NotEmpty(t, result.RefreshToken)
			assert.True(t, result.ExpiresAt.After(time.Now()))
		})
	}
}

func TestAuthService_LoginUser_InfrastructureErrorIsNotMasked(t *testing.T) {
	repo := &mockAuthRepository{
		getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
			return core_domain.User{}, errors.New("connection refused")
		},
	}

	service := newTestAuthService(t, repo)

	_, err := service.LoginUser(context.Background(), "user@example.com", "any-password")

	require.Error(t, err)
	assert.False(t, errors.Is(err, core_errors.ErrInvalidArgument),
		"инфраструктурная ошибка не должна маскироваться под неверные учётные данные")
}
