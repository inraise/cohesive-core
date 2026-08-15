package households_service

type HouseHoldsService struct {
	houseHoldsRepository HouseHoldsRepository
}

type HouseHoldsRepository interface {
}

func NewHouseHoldsService(
	houseHoldsRepository HouseHoldsRepository,
) *HouseHoldsService {
	return &HouseHoldsService{
		houseHoldsRepository: houseHoldsRepository,
	}
}
