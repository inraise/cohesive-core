package users_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	core_http_types "cohesive-core/internal/core/transport/http/types"
	"fmt"
	"net/http"
)

type PatchUserRequest struct {
	Email        core_http_types.Nullable[string] `json:"email"`
	PasswordHash core_http_types.Nullable[string] `json:"password_hash"`

	FirstName core_http_types.Nullable[string] `json:"first_name"`
	LastName  core_http_types.Nullable[string] `json:"last_name"`
	Age       core_http_types.Nullable[int]    `json:"age"`
}

func (r *PatchUserRequest) Validate() error {
	if r.Email.Set {
		if r.Email.Value == nil {
			return fmt.Errorf("`Email` can't be NULL")
		}

		emailLen := len([]rune(*r.Email.Value))

		if emailLen < 5 || emailLen > 100 {
			return fmt.Errorf("`Email` must be between 5 and 100 symbols")
		}
	}

	if r.PasswordHash.Set {
		if r.PasswordHash.Value == nil {
			return fmt.Errorf("`Password` can't be NULL")
		}
	}

	if r.FirstName.Set {
		if r.FirstName.Value == nil {
			return fmt.Errorf("`FirstName` can't be NULL")
		}

		firstNameLen := len([]rune(*r.FirstName.Value))

		if firstNameLen < 1 || firstNameLen > 100 {
			return fmt.Errorf("`FirstName` must be between 1 and 100 symbols")
		}
	}

	if r.LastName.Set {
		if r.LastName.Value == nil {
			return fmt.Errorf("`LastName` can't be NULL")
		}

		lastNameLen := len([]rune(*r.LastName.Value))

		if lastNameLen < 1 || lastNameLen > 100 {
			return fmt.Errorf("`LastName` must be between 1 and 100 symbols")
		}
	}

	if r.Age.Set {
		if r.Age.Value == nil {
			return fmt.Errorf("`Age` can't be NULL")
		}

		if *r.Age.Value < 0 || *r.Age.Value > 130 {
			return fmt.Errorf("`Age` must be between 0 and 130 symbols")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchMe(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchUserRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)
	userDomain, err := h.usersService.PatchMe(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) core_domain.UserPatch {
	return core_domain.NewUserPatch(
		request.Email.ToDomain(),
		request.PasswordHash.ToDomain(),
		request.FirstName.ToDomain(),
		request.LastName.ToDomain(),
		request.Age.ToDomain(),
	)
}
