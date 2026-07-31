package pkg_http_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

type PersonRole string

const (
	ApiVersion1 = "1"
)

type ApiVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	personRole PersonRole
}

func NewApiVersionRouter(apiVersion ApiVersion, personRole PersonRole) *ApiVersionRouter {
	return &ApiVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		personRole: personRole,
	}
}

func (r *ApiVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		prefix := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(prefix, route.Handler)
	}
}
