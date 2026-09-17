package auth_service_test

import (
	core_jwt "cohesive-core/internal/core/jwt"
	auth_service "cohesive-core/internal/features/auth/service"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
)

const testRefreshTTL = 720 * time.Hour

func newTestAuthService(t *testing.T, repo auth_service.AuthRepository) *auth_service.AuthService {
	t.Helper()

	tokenManager, err := core_jwt.NewTokenManager("test-secret-key-not-for-prod", 15*time.Minute)
	require.NoError(t, err)

	return auth_service.NewAuthService(repo, tokenManager, testRefreshTTL)
}
