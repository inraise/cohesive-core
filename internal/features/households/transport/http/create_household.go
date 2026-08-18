package households_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	households_service "cohesive-core/internal/features/households/service"
	"fmt"
	"net/http"
)

type CreateHouseholdResponse HouseholdDTOResponse

func (h *HouseholdsHTTPHandler) CreateHousehold(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke CreateHousehold handler")

	ownerID, ok := core_transport_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("user id not found in request context"),
			"internal server errors",
		)

		return
	}

	var request households_service.CreateHouseholdRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	householdDomain, err := h.householdsService.CreateHousehold(ctx, ownerID, request)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create household")

		return
	}

	response := CreateHouseholdResponse(householdDTOFromDomain(householdDomain, core_domain.HouseholdRoleOwner))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
