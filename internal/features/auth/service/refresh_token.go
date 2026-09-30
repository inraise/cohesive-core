package auth_service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"

	"github.com/google/uuid"
)

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}

func (s *AuthService) issueRefreshToken(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	plainToken, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	tokenDomain := core_domain.NewRefreshTokenUninitialized(
		userID,
		hashToken(plainToken),
		time.Now().Add(s.refreshTTL),
	)

	if _, err := s.authRepository.CreateRefreshToken(ctx, tokenDomain); err != nil {
		return "", fmt.Errorf("save refresh token: %w", err)
	}

	return plainToken, nil
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshTokenPlain string,
) (*LoginDTOResponse, error) {
	stored, err := s.authRepository.GetRefreshTokenByHash(ctx, hashToken(refreshTokenPlain))
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return nil, fmt.Errorf("invalid refresh token: %w", core_errors.ErrUnauthorized)
		}

		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if !stored.IsValid() {
		return nil, fmt.Errorf("refresh token is expired or revoked: %w", core_errors.ErrUnauthorized)
	}

	user, err := s.authRepository.GetUserByID(ctx, stored.UserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return nil, fmt.Errorf("user of refresh token not found: %w", core_errors.ErrUnauthorized)
		}

		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	if !user.IsVerified {
		return nil, fmt.Errorf("email is not verified: %w", core_errors.ErrForbidden)
	}

	if err = s.authRepository.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return nil, fmt.Errorf("revoke used refresh token: %w", err)
	}

	accessToken, expiresAt, err := s.tokenManager.GenerateAccessToken(stored.UserID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := s.issueRefreshToken(ctx, stored.UserID)
	if err != nil {
		return nil, fmt.Errorf("issue refresh token: %w", err)
	}

	return &LoginDTOResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
