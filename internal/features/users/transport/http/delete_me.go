package users_transport_http

import (
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

// DeleteMe godoc
// @Summary Удалить пользователя
// @Description Удалить текущего пользователя из системы
// @Tags users
// @Security ApiKeyAuth
// @Success 204 "Пользователь удалён"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /users/me [delete]
func (h *UsersHTTPHandler) DeleteMe(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_transport_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("user id not found in request context"),
			"internal server errors",
		)

		return
	}

	if err := h.usersService.DeleteMe(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)

		return
	}

	responseHandler.NoContentResponse()
}
