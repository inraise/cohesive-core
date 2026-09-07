package households_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ListInvitesResponse []InviteDTOResponse

// ListInvites godoc
// @Summary Список инвайтов
// @Description Получить список инвайтов дома, включая отозванные и исчерпанные. Доступно owner и admin
// @Tags households
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Success 200 {array} households_transport_http.InviteDTOResponse "Список инвайтов"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 403 {object} core_transport_http_response.ErrorResponse "Forbidden - role can't view invites"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found or not a member"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/invites [get]
func (h *HouseholdsHTTPHandler) ListInvites(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	callerID, ok := core_transport_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("user id not found in request context"), "internal server errors")

		return
	}

	householdID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse household id %q: %v: %w", r.PathValue("id"), err, core_errors.ErrInvalidArgument),
			"invalid household id",
		)

		return
	}

	invites, err := h.householdsService.ListInvites(ctx, householdID, callerID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list invites")

		return
	}

	response := make(ListInvitesResponse, 0, len(invites))
	for _, invite := range invites {
		response = append(response, inviteDTOFromDomain(invite))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
