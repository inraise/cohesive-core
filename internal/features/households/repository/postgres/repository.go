package households_repository_postgres

import core_pool "cohesive-core/internal/core/repository/postgres/pool"

type HouseHoldsRepository struct {
	pool core_pool.Pool
}

func NewHouseHoldsRepository(
	pool core_pool.Pool,
) *HouseHoldsRepository {
	return &HouseHoldsRepository{
		pool: pool,
	}
}
