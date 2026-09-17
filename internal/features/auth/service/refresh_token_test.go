package auth_service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_RefreshToken_Success(t *testing.T) {
	userID := uuid.New()
	storedTokenID := uuid.New()

	stored := core_domain.RefreshToken{
		ID:        storedTokenID,
		UserID:    userID,
		TokenHash: "irrelevant-in-this-test",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	var revokedID uuid.UUID
	var createdNewToken bool

	repo := &mockAuthRepository{
		getRefreshTokenByHashFunc: func(_ context.Context, _ string) (core_domain.RefreshToken, error) {
			return stored, nil
		},
		revokeRefreshTokenFunc: func(_ context.Context, id uuid.UUID) error {
			revokedID = id

			return nil
		},
		createRefreshTokenFunc: func(
			_ context.Context,
			token core_domain.RefreshToken,
		) (core_domain.RefreshToken, error) {
			createdNewToken = true
			token.ID = uuid.New()

			return token, nil
		},
	}

	service := newTestAuthService(t, repo)

	result, err := service.RefreshToken(context.Background(), "some-plain-refresh-token")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	assert.Equal(t, storedTokenID, revokedID, "старый refresh-токен должен быть отозван")
	assert.True(t, createdNewToken, "должен быть выпущен новый refresh-токен")
}

func TestAuthService_RefreshToken_Rejected(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name string
		repo *mockAuthRepository
	}{
		{
			name: "токен не найден",
			repo: &mockAuthRepository{
				getRefreshTokenByHashFunc: func(_ context.Context, _ string) (core_domain.RefreshToken, error) {
					return core_domain.RefreshToken{}, fmt.Errorf("token: %w", core_errors.ErrNotFound)
				},
			},
		},
		{
			name: "токен истёк",
			repo: &mockAuthRepository{
				getRefreshTokenByHashFunc: func(_ context.Context, _ string) (core_domain.RefreshToken, error) {
					return core_domain.RefreshToken{
						ID:        uuid.New(),
						UserID:    userID,
						ExpiresAt: time.Now().Add(-time.Hour),
					}, nil
				},
			},
		},
		{
			name: "токен отозван",
			repo: &mockAuthRepository{
				getRefreshTokenByHashFunc: func(_ context.Context, _ string) (core_domain.RefreshToken, error) {
					revokedAt := time.Now().Add(-time.Minute)

					return core_domain.RefreshToken{
						ID:        uuid.New(),
						UserID:    userID,
						ExpiresAt: time.Now().Add(time.Hour),
						RevokedAt: &revokedAt,
					}, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestAuthService(t, tt.repo)

			result, err := service.RefreshToken(context.Background(), "some-plain-refresh-token")

			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, core_errors.ErrUnauthorized)
		})
	}
}
