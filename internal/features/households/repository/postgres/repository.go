package households_repository_postgres

import core_pool "cohesive-core/internal/core/repository/postgres/pool"

type HouseholdsRepository struct {
	pool core_pool.Pool
}

func NewHouseholdsRepository(
	pool core_pool.Pool,
) *HouseholdsRepository {
	return &HouseholdsRepository{
		pool: pool,
	}
}
