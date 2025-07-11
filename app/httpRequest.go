package main

import (
	"strings"
)

type HttpRequest struct {
	Method  HttpMethod
	Target  string
	Headers string
	Body    string
}

func (request *HttpRequest) getHeader(title string) string {
	splitHeaders := strings.SplitSeq(request.Headers, "\r\n")
	for header := range splitHeaders {
		if after, ok := strings.CutPrefix(header, title+": "); ok {
			return after
		}
	}
	return ""
}

func parseRequest(reuqestBuffer []byte) *HttpRequest {
	requestStringParts := strings.Split(string(reuqestBuffer), "\r\n\r\n")
	requestLineAndHeaders := strings.Split(requestStringParts[0], "\r\n")
	requestLine := requestLineAndHeaders[0]
	requestTarget := strings.Split(requestLine, " ")[1]
	requestMethod := HttpMethod(strings.Split(requestLine, " ")[0])
	return &HttpRequest{
		Method:  requestMethod,
		Target:  requestTarget,
		Body:    requestStringParts[1],
		Headers: strings.Join(requestLineAndHeaders[1:], "\r\n"),
	}
}
