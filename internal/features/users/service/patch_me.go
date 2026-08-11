package users_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *UsersService) PatchMe(
	ctx context.Context,
	id uuid.UUID,
	patch core_domain.UserPatch,
) (core_domain.User, error) {
	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return core_domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchedUser, err := s.usersRepository.PatchMe(ctx, id, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchedUser, nil
}
