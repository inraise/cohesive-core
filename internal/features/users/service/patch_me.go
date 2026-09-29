package users_service

import (
	"context"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
	core_logger "cohesive-core/internal/core/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersService) PatchMe(
	ctx context.Context,
	id uuid.UUID,
	patch core_domain.UserPatch,
) (core_domain.User, error) {
	if patch.Password.Set && patch.Password.Value != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*patch.Password.Value), bcrypt.DefaultCost)
		if err != nil {
			return core_domain.User{}, fmt.Errorf("hash new password: %w", err)
		}

		hashed := string(hashedPassword)
		patch.Password.Value = &hashed
	}

	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	emailChanged := patch.Email.Set && patch.Email.Value != nil && *patch.Email.Value != user.Email

	if err = user.ApplyPatch(patch); err != nil {
		return core_domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	if emailChanged {
		user.IsVerified = false
	}

	patchedUser, err := s.usersRepository.PatchMe(ctx, id, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	// Письмо на новый адрес. Сбой не должен откатывать смену email: флаг is_verified уже false,
	// а письмо можно запросить повторно через /auth/resend-verification.
	if emailChanged {
		if err = s.emailVerifier.SendEmailVerification(ctx, patchedUser); err != nil {
			core_logger.FromContextOrNop(ctx).Error(
				"failed to send verification email after email change",
				zap.String("user_id", patchedUser.ID.String()),
				zap.Error(err),
			)
		}
	}

	return patchedUser, nil
}
