package auth_transport_http

import (
	"fmt"
	"net/http"

	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
)

type ResendVerificationRequest struct {
	Email string `json:"email"`
}

func (r *ResendVerificationRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("`Email` can't be NULL")
	}

	return nil
}

// ResendVerification godoc
// @Summary Повторная отправка письма с подтверждением
// @Description Отправляет новое письмо с ссылкой подтверждения. Всегда отвечает 204, не раскрывая, существует ли аккаунт и подтверждён ли он
// @Tags auth
// @Accept json
// @Param request body ResendVerificationRequest true "ResendVerification тело запроса"
// @Success 204 "Запрос принят"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 429 {object} core_transport_http_response.ErrorResponse "Too many requests"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /auth/resend-verification [post]
func (h *AuthHTTPHandler) ResendVerification(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	var request ResendVerificationRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	if err := h.authService.ResendVerification(ctx, request.Email); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to resend verification email",
		)

		return
	}

	responseHandler.NoContentResponse()
}
