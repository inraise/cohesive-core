package auth_service

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginDTOResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (s *AuthService) LoginUser(
	ctx context.Context,
	email, password string,
) (*LoginDTOResponse, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return nil, fmt.Errorf("invalid email or password: %w", core_errors.ErrInvalidArgument)
		}

		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, fmt.Errorf("invalid email or password: %w", core_errors.ErrInvalidArgument)
	}

	accessToken, expiresAt, err := s.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &LoginDTOResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}
