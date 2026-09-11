package shoppinglists_repository_postgres

import core_pool "cohesive-core/internal/core/repository/postgres/pool"

type ShoppingListsRepository struct {
	pool core_pool.Pool
}

func NewShoppingListsRepository(
	pool core_pool.Pool,
) *ShoppingListsRepository {
	return &ShoppingListsRepository{
		pool: pool,
	}
}
