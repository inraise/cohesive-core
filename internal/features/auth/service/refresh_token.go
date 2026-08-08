package auth_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func generateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}

func (s *AuthService) issueRefreshToken(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	plainToken, err := generateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	tokenDomain := core_domain.NewRefreshTokenUninitialized(
		userID,
		hashRefreshToken(plainToken),
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
	stored, err := s.authRepository.GetRefreshTokenByHash(ctx, hashRefreshToken(refreshTokenPlain))
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return nil, fmt.Errorf("invalid refresh token: %w", core_errors.ErrUnauthorized)
		}

		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if !stored.IsValid() {
		return nil, fmt.Errorf("refresh token is expired or revoked: %w", core_errors.ErrUnauthorized)
	}

	if err := s.authRepository.RevokeRefreshToken(ctx, stored.ID); err != nil {
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
