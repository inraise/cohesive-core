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

// RemoveMember godoc
// @Summary Убрать участника / выйти из дома
// @Description Удаляет участника из дома. Если user_id совпадает с вызывающим - выход из дома (владелец сначала должен передать владение). admin может убрать только member, owner - любого
// @Tags households
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Param user_id path string true "ID пользователя"
// @Success 204 "Участник удалён"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 403 {object} core_transport_http_response.ErrorResponse "Forbidden - insufficient role to remove this member"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household or member not found"
// @Failure 409 {object} core_transport_http_response.ErrorResponse "Sole owner must transfer ownership or delete household first"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/members/{user_id} [delete]
func (h *HouseholdsHTTPHandler) RemoveMember(rw http.ResponseWriter, r *http.Request) {
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

	targetID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse user id %q: %v: %w", r.PathValue("user_id"), err, core_errors.ErrInvalidArgument),
			"invalid user id",
		)

		return
	}

	if err := h.householdsService.RemoveMember(ctx, householdID, callerID, targetID); err != nil {
		responseHandler.ErrorResponse(err, "failed to remove member")

		return
	}

	responseHandler.NoContentResponse()
}
