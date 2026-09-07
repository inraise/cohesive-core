package households_transport_http

import (
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type ListHouseholdsResponse []HouseholdDTOResponse

// ListHouseholds godoc
// @Summary Список моих домов
// @Description Получить список домов, в которых состоит текущий пользователь, с ролью в каждом
// @Tags households
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} households_transport_http.HouseholdDTOResponse "Список домов"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households [get]
func (h *HouseholdsHTTPHandler) ListHouseholds(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke ListHouseholds handler")

	userID, ok := core_transport_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("user id not found in request context"),
			"internal server errors",
		)

		return
	}

	householdsWithRoles, err := h.householdsService.ListMyHouseholds(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list households")

		return
	}

	response := make(ListHouseholdsResponse, 0, len(householdsWithRoles))
	for _, hwr := range householdsWithRoles {
		response = append(response, householdDTOFromDomain(hwr.Household, hwr.Role))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
