package auth_transport_http

import (
	"fmt"
	"net/http"

	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
)

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

func (r *VerifyEmailRequest) Validate() error {
	if r.Token == "" {
		return fmt.Errorf("`Token` can't be NULL")
	}

	return nil
}

type VerifyEmailResponse UserDTOResponse

// VerifyEmail godoc
// @Summary Подтверждение email
// @Description Подтверждает email пользователя по токену из письма. Токен одноразовый и ограничен по времени
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "VerifyEmail тело запроса"
// @Success 200 {object} VerifyEmailResponse "Email подтверждён"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Токен невалиден, просрочен или уже использован"
// @Failure 429 {object} core_transport_http_response.ErrorResponse "Too many requests"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /auth/verify-email [post]
func (h *AuthHTTPHandler) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	var request VerifyEmailRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	user, err := h.authService.VerifyEmail(ctx, request.Token)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to verify email",
		)

		return
	}

	responseHandler.JSONResponse(VerifyEmailResponse(userDTOFromDomain(user)), http.StatusOK)
}
