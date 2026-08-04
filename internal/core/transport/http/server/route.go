package core_transport_http_server

import (
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	"net/http"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_transport_http_middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_transport_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}
