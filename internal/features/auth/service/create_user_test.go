package auth_service_test

import (
	core_domain "cohesive-core/internal/core/domain"
	auth_service "cohesive-core/internal/features/auth/service"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_CreateUser(t *testing.T) {
	t.Run("пароль хешируется перед сохранением", func(t *testing.T) {
		const plainPassword = "supersecurepassword"

		var savedUser core_domain.User

		repo := &mockAuthRepository{
			createUserFunc: func(_ context.Context, user core_domain.User) (core_domain.User, error) {
				savedUser = user
				user.ID = uuid.New()

				return user, nil
			},
		}

		service := newTestAuthService(t, repo)

		request := auth_service.CreateUserRequest{
			Email:     "user@example.com",
			Password:  plainPassword,
			FirstName: "Nikita",
		}

		result, err := service.CreateUser(context.Background(), request)

		require.NoError(t, err)

		assert.NotEqual(t, plainPassword, savedUser.PasswordHash)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(savedUser.PasswordHash), []byte(plainPassword)))
		assert.Equal(t, "user@example.com", result.Email)
	})

	t.Run("невалидные данные отклоняются до похода в репозиторий", func(t *testing.T) {
		repo := &mockAuthRepository{
			createUserFunc: func(_ context.Context, _ core_domain.User) (core_domain.User, error) {
				t.Fatal("CreateUser репозитория не должен вызываться для заведомо невалидных данных")

				return core_domain.User{}, nil
			},
		}

		service := newTestAuthService(t, repo)

		request := auth_service.CreateUserRequest{
			Email:     "a",
			Password:  "supersecurepassword",
			FirstName: "Nikita",
		}

		_, err := service.CreateUser(context.Background(), request)

		require.Error(t, err)
	})
}
