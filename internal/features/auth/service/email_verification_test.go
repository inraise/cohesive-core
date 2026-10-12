package auth_service_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	auth_service "cohesive-core/internal/features/auth/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var verificationLink = regexp.MustCompile(`https://app\.example\.com/verify-email\?token=([A-Za-z0-9_-]+)`)

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sum[:])
}

func extractToken(t *testing.T, body string) string {
	t.Helper()

	match := verificationLink.FindStringSubmatch(body)
	require.Len(t, match, 2, "в теле письма нет ссылки подтверждения: %q", body)

	return match[1]
}

func newUnverifiedUser(email string) core_domain.User {
	return core_domain.NewUser(
		uuid.New(),
		1,
		email,
		"password-hash",
		"Nikita",
		nil,
		nil,
		false,
		time.Now(),
		time.Now(),
	)
}

func newRegisterRequest() auth_service.CreateUserRequest {
	return auth_service.CreateUserRequest{
		Email:     "user@example.com",
		Password:  "supersecurepassword",
		FirstName: "Nikita",
	}
}

func TestAuthService_CreateUser_SendsVerificationEmail(t *testing.T) {
	var (
		createdUserID = uuid.New()
		saved         core_domain.EmailVerification
	)

	repo := &mockAuthRepository{
		createUserFunc: func(_ context.Context, user core_domain.User) (core_domain.User, error) {
			user.ID = createdUserID

			return user, nil
		},
		createEmailVerificationFunc: func(
			_ context.Context,
			verification core_domain.EmailVerification,
		) (core_domain.EmailVerification, error) {
			saved = verification

			return verification, nil
		},
	}
	mailer := newMockMailer()
	service := newTestAuthServiceWithMailer(t, repo, mailer)

	_, err := service.CreateUser(context.Background(), newRegisterRequest())
	require.NoError(t, err)

	email := waitForEmail(t, mailer)
	assert.Equal(t, "user@example.com", email.to)
	assert.NotEmpty(t, email.subject)

	token := extractToken(t, email.body)

	assert.Equal(t, createdUserID, saved.UserID)
	assert.Equal(t, "user@example.com", saved.Email)
	assert.Equal(t, sha256Hex(token), saved.TokenHash, "в БД должен лежать SHA-256 токена из письма")
	assert.NotEqual(t, token, saved.TokenHash, "открытый токен не должен попадать в БД")
	assert.WithinDuration(t, time.Now().Add(testVerificationTTL), saved.ExpiresAt, time.Minute)
}

func TestAuthService_CreateUser_VerificationFailureDoesNotFailRegistration(t *testing.T) {
	t.Run("сбой записи верификации в БД", func(t *testing.T) {
		userID := uuid.New()

		repo := &mockAuthRepository{
			createUserFunc: func(_ context.Context, user core_domain.User) (core_domain.User, error) {
				user.ID = userID

				return user, nil
			},
			createEmailVerificationFunc: func(
				_ context.Context,
				_ core_domain.EmailVerification,
			) (core_domain.EmailVerification, error) {
				return core_domain.EmailVerification{}, errors.New("db is down")
			},
		}
		mailer := newMockMailer()
		service := newTestAuthServiceWithMailer(t, repo, mailer)

		result, err := service.CreateUser(context.Background(), newRegisterRequest())

		require.NoError(t, err)
		assert.Equal(t, userID, result.ID)
		assertNoEmail(t, mailer)
	})

	t.Run("сбой отправки письма", func(t *testing.T) {
		userID := uuid.New()

		repo := &mockAuthRepository{
			createUserFunc: func(_ context.Context, user core_domain.User) (core_domain.User, error) {
				user.ID = userID

				return user, nil
			},
			createEmailVerificationFunc: func(
				_ context.Context,
				verification core_domain.EmailVerification,
			) (core_domain.EmailVerification, error) {
				return verification, nil
			},
		}
		mailer := newMockMailer()
		mailer.err = errors.New("smtp is down")
		service := newTestAuthServiceWithMailer(t, repo, mailer)

		result, err := service.CreateUser(context.Background(), newRegisterRequest())

		require.NoError(t, err)
		assert.Equal(t, userID, result.ID)
		waitForEmail(t, mailer)
	})
}

func TestAuthService_VerifyEmail(t *testing.T) {
	const plainToken = "token-from-email"

	t.Run("репозиторию уходит хеш токена, а не сам токен", func(t *testing.T) {
		var receivedHash string

		verifiedUser := newUnverifiedUser("user@example.com")
		verifiedUser.IsVerified = true

		repo := &mockAuthRepository{
			verifyEmailByHashFunc: func(_ context.Context, tokenHash string) (core_domain.User, error) {
				receivedHash = tokenHash

				return verifiedUser, nil
			},
		}

		service := newTestAuthService(t, repo)

		result, err := service.VerifyEmail(context.Background(), plainToken)

		require.NoError(t, err)
		assert.Equal(t, sha256Hex(plainToken), receivedHash)
		assert.True(t, result.IsVerified)
	})

	t.Run("невалидный токен возвращает ErrInvalidArgument", func(t *testing.T) {
		repo := &mockAuthRepository{
			verifyEmailByHashFunc: func(_ context.Context, _ string) (core_domain.User, error) {
				return core_domain.User{}, fmt.Errorf("token is invalid: %w", core_errors.ErrInvalidArgument)
			},
		}

		service := newTestAuthService(t, repo)

		_, err := service.VerifyEmail(context.Background(), plainToken)

		require.Error(t, err)
		assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})
}

func TestAuthService_ResendVerification(t *testing.T) {
	t.Run("неподтверждённому пользователю уходит новое письмо", func(t *testing.T) {
		user := newUnverifiedUser("user@example.com")

		var saved core_domain.EmailVerification

		repo := &mockAuthRepository{
			getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
				return user, nil
			},
			createEmailVerificationFunc: func(
				_ context.Context,
				verification core_domain.EmailVerification,
			) (core_domain.EmailVerification, error) {
				saved = verification

				return verification, nil
			},
		}
		mailer := newMockMailer()
		service := newTestAuthServiceWithMailer(t, repo, mailer)

		require.NoError(t, service.ResendVerification(context.Background(), user.Email))

		email := waitForEmail(t, mailer)
		assert.Equal(t, user.Email, email.to)
		assert.Equal(t, user.ID, saved.UserID)
		assert.Equal(t, sha256Hex(extractToken(t, email.body)), saved.TokenHash)
	})

	t.Run("уже подтверждённому пользователю письмо не отправляется", func(t *testing.T) {
		user := newUnverifiedUser("user@example.com")
		user.IsVerified = true

		repo := &mockAuthRepository{
			getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
				return user, nil
			},
		}
		mailer := newMockMailer()
		service := newTestAuthServiceWithMailer(t, repo, mailer)

		require.NoError(t, service.ResendVerification(context.Background(), user.Email))
		assertNoEmail(t, mailer)
	})

	t.Run("несуществующий email не раскрывается: nil и без письма", func(t *testing.T) {
		repo := &mockAuthRepository{
			getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
				return core_domain.User{}, fmt.Errorf("user lookup: %w", core_errors.ErrNotFound)
			},
		}
		mailer := newMockMailer()
		service := newTestAuthServiceWithMailer(t, repo, mailer)

		require.NoError(t, service.ResendVerification(context.Background(), "ghost@example.com"))
		assertNoEmail(t, mailer)
	})

	t.Run("инфраструктурная ошибка не скрывается", func(t *testing.T) {
		repo := &mockAuthRepository{
			getUserByEmailFunc: func(_ context.Context, _ string) (core_domain.User, error) {
				return core_domain.User{}, errors.New("connection refused")
			},
		}
		service := newTestAuthService(t, repo)

		err := service.ResendVerification(context.Background(), "user@example.com")

		require.Error(t, err)
		assert.False(t, errors.Is(err, core_errors.ErrNotFound))
	})
}
