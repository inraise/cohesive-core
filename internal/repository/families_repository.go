package repository

import (
	"cohesive-core/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (r *FamilyRepository) JoinFamily(
	ctx context.Context,
	inviteCode string,
	userID string,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var familyID uuid.UUID
	familyQuery := `SELECT id FROM families WHERE invite_code = $1`
	err = tx.QueryRow(ctx, familyQuery, inviteCode).Scan(&familyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("Неверный инвайт-код")
		}
		return err
	}

	memberQuery := `
		INSERT INTO family_members (family_id, user_id, role, joined_at)
		VALUES ($1, $2, 'member', NOW())
		ON CONFLICT (family_id, user_id) DO NOTHING`

	res, err := tx.Exec(ctx, memberQuery, familyID, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("Вы уже состоите в этой семье")
	}

	return tx.Commit(ctx)
}

func (r *FamilyRepository) LeaveFamily(ctx context.Context, userID string) error {
	query := `DELETE FROM family_members WHERE user_id = $1`
	res, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("Вы не состоите ни в одной семье")
	}
	return nil
}

func (r *FamilyRepository) UpdateMemberRole(
	ctx context.Context,
	familyID, targetUserID string,
	newRole string,
) error {
	query := `
		UPDATE family_members 
		SET role = $1::role_members 
		WHERE family_id = $2 AND user_id = $3`

	res, err := r.db.Exec(ctx, query, newRole, familyID, targetUserID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("Участник не найден в этой семье")
	}
	return nil
}

func (r *FamilyRepository) KickMember(ctx context.Context, familyID, targetUserID string) error {
	query := `DELETE FROM family_members WHERE family_id = $1 AND user_id = $2`
	res, err := r.db.Exec(ctx, query, familyID, targetUserID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("Участник не найден в этой семье")
	}
	return nil
}

func (r *FamilyRepository) IsAdmin(ctx context.Context, familyID, userID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM family_members 
			WHERE family_id = $1 AND user_id = $2 AND role = 'admin'
		)`
	var isAdmin bool
	err := r.db.QueryRow(ctx, query, familyID, userID).Scan(&isAdmin)
	return isAdmin, err
}
