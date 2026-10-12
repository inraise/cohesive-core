package auth_service_test

import (
	"testing"
	"time"

	core_jwt "cohesive-core/internal/core/jwt"
	core_mailer "cohesive-core/internal/core/mailer"
	core_verification "cohesive-core/internal/core/verification"
	auth_service "cohesive-core/internal/features/auth/service"

	"github.com/go-openapi/testify/v2/require"
)

const (
	testRefreshTTL      = 720 * time.Hour
	testVerificationURL = "https://app.example.com/verify-email"
	testVerificationTTL = 24 * time.Hour
)

func newTestAuthService(t *testing.T, repo auth_service.AuthRepository) *auth_service.AuthService {
	t.Helper()

	return newTestAuthServiceWithMailer(t, repo, newMockMailer())
}

func newTestAuthServiceWithMailer(
	t *testing.T,
	repo auth_service.AuthRepository,
	mailer core_mailer.Mailer,
) *auth_service.AuthService {
	t.Helper()

	tokenManager, err := core_jwt.NewTokenManager("test-secret-key-not-for-prod", 15*time.Minute)
	require.NoError(t, err)

	return auth_service.NewAuthService(
		repo,
		tokenManager,
		testRefreshTTL,
		mailer,
		core_verification.Config{
			BaseURL:  testVerificationURL,
			TokenTTL: testVerificationTTL,
		},
	)
}
