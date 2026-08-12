package users_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *UsersService) GetMe(
	ctx context.Context,
	userID uuid.UUID,
) (core_domain.User, error) {
	user, err := s.usersRepository.GetUserByID(ctx, userID)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
