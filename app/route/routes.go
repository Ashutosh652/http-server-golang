package route

import (
	"http-server-go/app/config"
	"http-server-go/app/http"
	"net"
)

type Route struct {
	Path           string
	MatchExact     bool
	Handler        func(conn net.Conn, config config.Config, req *http.HttpRequest) *http.HttpResponse
	AllowedMethods http.HttpMethodList
}

func InitRoutes() []Route {
	routes := []Route{
		{"/", true, handleRoot, []http.HttpMethod{http.HTTP_GET_METHOD}},
		{"/user-agent", true, handleUserAgent, []http.HttpMethod{http.HTTP_GET_METHOD}},
		{"/echo", false, handleEcho, []http.HttpMethod{http.HTTP_GET_METHOD}},
		{"/files", false, handleFiles, []http.HttpMethod{http.HTTP_GET_METHOD, http.HTTP_POST_METHOD}},
	}
	return routes
}
