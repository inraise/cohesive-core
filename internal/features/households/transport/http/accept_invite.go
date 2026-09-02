package households_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

func (h *HouseholdsHTTPHandler) AcceptInvite(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_transport_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("user id not found in request context"), "internal server errors")

		return
	}

	code := r.PathValue("code")

	household, err := h.householdsService.AcceptInvite(ctx, code, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to accept invite")

		return
	}

	response := householdDTOFromDomain(household, core_domain.HouseholdRoleMember)
	responseHandler.JSONResponse(response, http.StatusOK)
}
