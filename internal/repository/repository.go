package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	Auth  *AuthRepository
	Family *FamilyRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Auth:  NewAuthRepository(db),
		Family: NewFamilyRepository(db),
	}
}
