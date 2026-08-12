package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *UsersService) DeleteMe(
	ctx context.Context,
	id uuid.UUID,
) error {
	if err := s.usersRepository.DeleteMe(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
