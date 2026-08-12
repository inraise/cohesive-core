package users_repository_postgres

import core_pool "cohesive-core/internal/core/repository/postgres/pool"

type UsersRepository struct {
	pool core_pool.Pool
}

func NewUsersRepository(
	pool core_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
