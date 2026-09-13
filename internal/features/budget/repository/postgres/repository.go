package budget_repository_postgres

import core_pool "cohesive-core/internal/core/repository/postgres/pool"

type BudgetRepository struct {
	pool core_pool.Pool
}

func NewBudgetRepository(
	pool core_pool.Pool,
) *BudgetRepository {
	return &BudgetRepository{
		pool: pool,
	}
}
