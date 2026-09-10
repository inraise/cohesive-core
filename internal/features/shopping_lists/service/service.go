package shoppinglists_service

type ShoppingListsService struct {
	shoppingListsRepository ShoppingListsRepository
}

type ShoppingListsRepository interface {
}

func NewShoppingListsService(
	shoppingListsRepository ShoppingListsRepository,
) *ShoppingListsService {
	return &ShoppingListsService{
		shoppingListsRepository: shoppingListsRepository,
	}
}
