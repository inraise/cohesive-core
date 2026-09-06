package tasks_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	core_http_types "cohesive-core/internal/core/transport/http/types"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string]    `json:"title"`
	Description core_http_types.Nullable[string]    `json:"description"`
	Status      core_http_types.Nullable[string]    `json:"status"`
	AssignedTo  core_http_types.Nullable[uuid.UUID] `json:"assigned_to"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 200 {
			return fmt.Errorf("`Title` must be between 1 and 200 symbols")
		}
	}

	if r.Status.Set {
		if r.Status.Value == nil {
			return fmt.Errorf("`Status` can't be NULL")
		}

		switch core_domain.TaskStatus(*r.Status.Value) {
		case core_domain.TaskStatusPending, core_domain.TaskStatusInProgress, core_domain.TaskStatusDone:
		default:
			return fmt.Errorf("`Status` must be one of: pending, in_progress, done")
		}
	}

	return nil
}

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
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

	taskID, err := uuid.Parse(r.PathValue("task_id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse task id %q: %v: %w", r.PathValue("task_id"), err, core_errors.ErrInvalidArgument),
			"invalid task id",
		)

		return
	}

	var request PatchTaskRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	patch := core_domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Status.ToDomain(),
		request.AssignedTo.ToDomain(),
	)

	task, err := h.tasksService.PatchTask(ctx, householdID, taskID, callerID, patch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")

		return
	}

	response := taskDTOFromDomain(task)
	responseHandler.JSONResponse(response, http.StatusOK)
}
