package users_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	users_service "cohesive-core/internal/features/users/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUsersRepository struct {
	stored core_domain.User
}

func (m *mockUsersRepository) GetUserByID(_ context.Context, _ uuid.UUID) (core_domain.User, error) {
	return m.stored, nil
}

func (m *mockUsersRepository) PatchMe(
	_ context.Context,
	_ uuid.UUID,
	user core_domain.User,
) (core_domain.User, error) {
	return user, nil
}

func (m *mockUsersRepository) DeleteMe(_ context.Context, _ uuid.UUID) error {
	panic("DeleteMe is not expected in these tests")
}

type mockEmailVerifier struct {
	calls []core_domain.User
	err   error
}

func (m *mockEmailVerifier) SendEmailVerification(_ context.Context, user core_domain.User) error {
	m.calls = append(m.calls, user)

	return m.err
}

func newStoredUser() core_domain.User {
	return core_domain.NewUser(
		uuid.New(),
		1,
		"old@example.com",
		"password-hash",
		"Nikita",
		nil,
		nil,
		true,
		time.Now().Add(-time.Hour),
		time.Now().Add(-time.Hour),
	)
}

func emailPatch(email string) core_domain.UserPatch {
	return core_domain.UserPatch{
		Email: core_domain.Nullable[string]{Set: true, Value: &email},
	}
}

func TestUsersService_PatchMe_EmailChange(t *testing.T) {
	t.Run("смена email сбрасывает is_verified и отправляет письмо на новый адрес", func(t *testing.T) {
		stored := newStoredUser()
		verifier := &mockEmailVerifier{}
		service := users_service.NewUsersService(&mockUsersRepository{stored: stored}, verifier)

		result, err := service.PatchMe(context.Background(), stored.ID, emailPatch("new@example.com"))

		require.NoError(t, err)
		assert.False(t, result.IsVerified)
		require.Len(t, verifier.calls, 1)
		assert.Equal(t, "new@example.com", verifier.calls[0].Email)
		assert.Equal(t, stored.ID, verifier.calls[0].ID)
	})

	t.Run("тот же email не сбрасывает подтверждение и не шлёт письмо", func(t *testing.T) {
		stored := newStoredUser()
		verifier := &mockEmailVerifier{}
		service := users_service.NewUsersService(&mockUsersRepository{stored: stored}, verifier)

		result, err := service.PatchMe(context.Background(), stored.ID, emailPatch(stored.Email))

		require.NoError(t, err)
		assert.True(t, result.IsVerified)
		assert.Empty(t, verifier.calls)
	})

	t.Run("правка имени письма не вызывает", func(t *testing.T) {
		stored := newStoredUser()
		verifier := &mockEmailVerifier{}
		service := users_service.NewUsersService(&mockUsersRepository{stored: stored}, verifier)

		firstName := "Daniil"
		patch := core_domain.UserPatch{
			FirstName: core_domain.Nullable[string]{Set: true, Value: &firstName},
		}

		result, err := service.PatchMe(context.Background(), stored.ID, patch)

		require.NoError(t, err)
		assert.True(t, result.IsVerified)
		assert.Empty(t, verifier.calls)
	})

	t.Run("сбой отправки письма не откатывает смену email", func(t *testing.T) {
		stored := newStoredUser()
		verifier := &mockEmailVerifier{err: errors.New("smtp is down")}
		service := users_service.NewUsersService(&mockUsersRepository{stored: stored}, verifier)

		result, err := service.PatchMe(context.Background(), stored.ID, emailPatch("new@example.com"))

		require.NoError(t, err)
		assert.Equal(t, "new@example.com", result.Email)
		assert.False(t, result.IsVerified)
	})
}
