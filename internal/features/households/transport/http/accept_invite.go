package households_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

// AcceptInvite godoc
// @Summary Принять инвайт
// @Description Присоединиться к дому по коду приглашения, становится member. Не требует предварительного членства в доме - только валидный access-токен
// @Tags households
// @Produce json
// @Security ApiKeyAuth
// @Param code path string true "Код приглашения"
// @Success 200 {object} households_transport_http.HouseholdDTOResponse "Приглашение принято"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Invite is invalid, expired, exhausted, or already accepted"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/invites/{code}/accept [post]
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
