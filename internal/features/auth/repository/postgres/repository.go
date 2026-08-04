package auth_repository_postgres

import core_pool "cohesive-core/internal/core/repository/postgres/pool"

type AuthRepository struct {
	pool core_pool.Pool
}

func NewAuthRepository(
	pool core_pool.Pool,
) *AuthRepository {
	return &AuthRepository{
		pool: pool,
	}
}
