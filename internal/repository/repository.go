package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	Auth  *AuthRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Auth:  NewAuthRepository(db),
	}
}
