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

// RevokeInvite godoc
// @Summary Отозвать инвайт
// @Description Отозвать инвайт-код, дальше по нему нельзя присоединиться. Доступно owner и admin
// @Tags households
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Param invite_id path string true "ID инвайта"
// @Success 204 "Инвайт отозван"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 403 {object} core_transport_http_response.ErrorResponse "Forbidden - role can't revoke invites"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household, member, or invite not found"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/invites/{invite_id} [delete]
func (h *HouseholdsHTTPHandler) RevokeInvite(rw http.ResponseWriter, r *http.Request) {
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

	inviteID, err := uuid.Parse(r.PathValue("invite_id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse invite id %q: %v: %w", r.PathValue("invite_id"), err, core_errors.ErrInvalidArgument),
			"invalid invite id",
		)

		return
	}

	if err := h.householdsService.RevokeInvite(ctx, householdID, callerID, inviteID); err != nil {
		responseHandler.ErrorResponse(err, "failed to revoke invite")

		return
	}

	responseHandler.NoContentResponse()
}
