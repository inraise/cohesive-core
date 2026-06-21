package service

import (
	"cohesive-core/internal/models"
	"cohesive-core/internal/repository"
	"context"
	"crypto/rand"

	"github.com/google/uuid"
)

type FamilyService struct {
	repo *repository.FamilyRepository
}

func NewFamilyService(repo *repository.FamilyRepository) *FamilyService {
	return &FamilyService{
		repo: repo,
	}
}

func (s *FamilyService) CreateFamily(
	ctx context.Context,
	req models.CreateFamilyRequest,
	userID uuid.UUID,
) (*models.Family, error) {
	inviteCode, err := generateInviteCode(8)
	if err != nil {
		return nil, err
	}

	family := &models.Family{
		Name:       req.Name,
		InviteCode: inviteCode,
	}

	err = s.repo.CreateFamily(ctx, family, userID)
	if err != nil {
		return nil, err
	}

	return family, nil
}

func generateInviteCode(n int) (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
