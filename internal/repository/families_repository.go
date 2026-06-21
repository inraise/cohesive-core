package repository

import (
	"cohesive-core/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
POST /api/v1/family — Создать новую семью
PUT /api/v1/family — Обновить данные семьи (название)
*/

type FamilyRepository struct {
	db *pgxpool.Pool
}

func NewFamilyRepository(db *pgxpool.Pool) *FamilyRepository {
	return &FamilyRepository{
		db: db,
	}
}

func (r *FamilyRepository) CreateFamily(
	ctx context.Context,
	family *models.Family,
	userId uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	familyQuery := `
		INSERT INTO families (name, invite_code, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, familyQuery, family.Name, family.InviteCode).
		Scan(&family.Id, &family.CreatedAt, &family.UpdatedAt)
	if err != nil {
		return err
	}

	memberQuery := `
		INSERT INTO family_members (family_id, user_id, role, joined_at)
		VALUES ($1, $2, 'admin', NOW())`

	_, err = tx.Exec(ctx, memberQuery, family.Id, userId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *FamilyRepository) UpdateFamilyName(
	ctx context.Context,
	familyID uuid.UUID,
	userID uuid.UUID,
	newName string,
) error {
	query := `
		UPDATE families 
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND EXISTS (
			SELECT 1 FROM family_members 
			WHERE family_id = $2 AND user_id = $3 AND role = 'admin'
		)`

	res, err := r.db.Exec(ctx, query, newName, familyID, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
