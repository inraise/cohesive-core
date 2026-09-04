package tasks_transport_http

import core_transport_http_server "cohesive-core/internal/core/transport/http/server"

type TasksHTTPHandler struct {
	tasksService TaskService
}

type TaskService interface {
}

func (h *TasksHTTPHandler) Routes() []core_transport_http_server.Route {
	return []core_transport_http_server.Route{}
}
