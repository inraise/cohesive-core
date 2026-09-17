package auth_service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_LogoutUser_RevokesToken(t *testing.T) {
	storedTokenID := uuid.New()
	var revokedID uuid.UUID

	repo := &mockAuthRepository{
		getRefreshTokenByHashFunc: func(_ context.Context, _ string) (core_domain.RefreshToken, error) {
			return core_domain.RefreshToken{ID: storedTokenID}, nil
		},
		revokeRefreshTokenFunc: func(_ context.Context, id uuid.UUID) error {
			revokedID = id

			return nil
		},
	}

	service := newTestAuthService(t, repo)

	err := service.LogoutUser(context.Background(), "some-plain-refresh-token")

	require.NoError(t, err)
	assert.Equal(t, storedTokenID, revokedID)
}

func TestAuthService_LogoutUser_MissingTokenIsNotAnError(t *testing.T) {
	repo := &mockAuthRepository{
		getRefreshTokenByHashFunc: func(_ context.Context, _ string) (core_domain.RefreshToken, error) {
			return core_domain.RefreshToken{}, fmt.Errorf("token: %w", core_errors.ErrNotFound)
		},
		revokeRefreshTokenFunc: func(_ context.Context, _ uuid.UUID) error {
			t.Fatal("RevokeRefreshToken не должен вызываться, если токен и так не найден")

			return nil
		},
	}

	service := newTestAuthService(t, repo)

	err := service.LogoutUser(context.Background(), "already-gone-token")

	require.NoError(t, err)
}

func TestAuthService_LogoutUser_InfrastructureErrorPropagates(t *testing.T) {
	repo := &mockAuthRepository{
		getRefreshTokenByHashFunc: func(_ context.Context, _ string) (core_domain.RefreshToken, error) {
			return core_domain.RefreshToken{}, errors.New("connection refused")
		},
	}

	service := newTestAuthService(t, repo)

	err := service.LogoutUser(context.Background(), "some-token")

	require.Error(t, err)
}
