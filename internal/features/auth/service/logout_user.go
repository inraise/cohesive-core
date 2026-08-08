package auth_service

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"errors"
	"fmt"
)

func (s *AuthService) LogoutUser(
	ctx context.Context,
	refreshTokenPlain string,
) error {
	stored, err := s.authRepository.GetRefreshTokenByHash(ctx, hashRefreshToken(refreshTokenPlain))
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("get refresh token: %w", err)
	}

	if err := s.authRepository.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
