package core_transport_http_server

import (
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middleware []core_transport_http_middleware.Middleware
}

func NewAPIVersionRouter(
	apiVersion ApiVersion,
	middleware ...core_transport_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_transport_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
