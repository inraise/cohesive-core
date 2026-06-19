package repository

import (
	"cohesive-core/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
Модуль 1: Аутентификация и Управление Семьей
• POST /api/v1/auth/register — Регистрация
• POST /api/v1/auth/login — Логин
• POST /api/v1/family — Создать новую семью
• PUT /api/v1/family — Обновить данные семьи (название)
• POST /api/v1/family/join — Войти по invite_code
• DELETE /api/v1/family/leave — Выйти из семьи
• PATCH /api/v1/family/members/{user_id} — Изменить роль (admin/member/child)
• DELETE /api/v1/family/members/{user_id} — Удалить человека из семьи
*/

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, created_at, updated_at
		FROM users 
		WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
