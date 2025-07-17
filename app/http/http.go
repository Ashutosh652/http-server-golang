package http

import "strings"

type HttpMethod string
type HttpMethodList []HttpMethod

const (
	HTTP_GET_METHOD  HttpMethod = "GET"
	HTTP_POST_METHOD HttpMethod = "POST"
)

type BaseHttp struct {
	Headers string
	Body    string
}

func (methods HttpMethodList) AsStrings() []string {
	strs := make([]string, len(methods))
	for i, m := range methods {
		strs[i] = string(m)
	}
	return strs
}

func (http *BaseHttp) GetHeader(title string) string {
	splitHeaders := strings.SplitSeq(http.Headers, "\r\n")
	for header := range splitHeaders {
		if after, ok := strings.CutPrefix(header, title+": "); ok {
			return after
		}
	}
	return ""
}
