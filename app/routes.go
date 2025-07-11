package main

import "net"

type HttpMethod string
type HttpMethodList []HttpMethod

const (
	HTTP_GET_METHOD  HttpMethod = "GET"
	HTTP_POST_METHOD HttpMethod = "POST"
)

type Route struct {
	Path           string
	MatchExact     bool
	Handler        func(conn net.Conn, config Config, req *HttpRequest) *HttpResponse
	AllowedMethods HttpMethodList
}

func (methods HttpMethodList) AsStrings() []string {
	strs := make([]string, len(methods))
	for i, m := range methods {
		strs[i] = string(m)
	}
	return strs
}

func initRoutes() []Route {
	routes := []Route{
		{"/", true, handleRoot, []HttpMethod{HTTP_GET_METHOD}},
		{"/user-agent", true, handleUserAgent, []HttpMethod{HTTP_GET_METHOD}},
		{"/echo", false, handleEcho, []HttpMethod{HTTP_GET_METHOD}},
		{"/files", false, handleFiles, []HttpMethod{HTTP_GET_METHOD, HTTP_POST_METHOD}},
	}
	return routes
}
