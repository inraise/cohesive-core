package households_transport_http

import (
	core_jwt "cohesive-core/internal/core/jwt"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
)

type HouseHoldsHTTPHandler struct {
	houseHoldsService HouseHoldsService
	tokenManager      *core_jwt.TokenManager
}

type HouseHoldsService interface {
}

func NewHouseHoldsHTTPHandler(
	houseHoldsService HouseHoldsService,
	tokenManager *core_jwt.TokenManager,
) *HouseHoldsHTTPHandler {
	return &HouseHoldsHTTPHandler{
		houseHoldsService: houseHoldsService,
		tokenManager:      tokenManager,
	}
}

func (h *HouseHoldsHTTPHandler) Routes() []core_transport_http_server.Route {
	// authenticate := core_transport_http_middleware.Authenticate(h.tokenManager)

	return []core_transport_http_server.Route{}
}
