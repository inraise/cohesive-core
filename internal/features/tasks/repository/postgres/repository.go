package tasks_repository_postgres

import core_pool "cohesive-core/internal/core/repository/postgres/pool"

type TasksRepository struct {
	pool core_pool.Pool
}

func NewAuthRepository(
	pool core_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
