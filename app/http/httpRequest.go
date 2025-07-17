package http

import (
	"strings"
)

type HttpRequest struct {
	Method HttpMethod
	Target string
	BaseHttp
}

func ParseRequest(reuqestBuffer []byte) *HttpRequest {
	requestStringParts := strings.Split(string(reuqestBuffer), "\r\n\r\n")
	requestLineAndHeaders := strings.Split(requestStringParts[0], "\r\n")
	requestLine := requestLineAndHeaders[0]
	requestTarget := strings.Split(requestLine, " ")[1]
	requestMethod := HttpMethod(strings.Split(requestLine, " ")[0])
	return &HttpRequest{
		Method: requestMethod,
		Target: requestTarget,
		BaseHttp: BaseHttp{
			Body:    requestStringParts[1],
			Headers: strings.Join(requestLineAndHeaders[1:], "\r\n"),
		},
	}
}
