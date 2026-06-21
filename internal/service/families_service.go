package service

import (
	"cohesive-core/internal/models"
	"cohesive-core/internal/repository"
	"context"
	"crypto/rand"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	req models.FamilyRequest,
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

func (s *FamilyService) UpdateFamily(
	ctx context.Context,
	familyID uuid.UUID,
	userID uuid.UUID,
	req models.FamilyRequest,
) error {
	if req.Name == "" {
		return errors.New("Название семьи не может быть пустым")
	}

	err := s.repo.UpdateFamilyName(ctx, familyID, userID, req.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("У вас нет прав на редактирование этой семьи или она не существует")
		}
		return err
	}

	return nil
}
