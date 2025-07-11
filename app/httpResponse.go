package main

import (
	"fmt"
	"strconv"
	"strings"
)

type HttpResponse struct {
	Version      string
	Status       int
	ReasonPhrase string
	Headers      string
	Body         string
}

func (httpResponse *HttpResponse) createResponse() string {
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s", httpResponse.Version, httpResponse.Status, httpResponse.ReasonPhrase, httpResponse.Headers, httpResponse.Body)
}

func (httpResponse *HttpResponse) headerExists(title string) bool {
	splitHeaders := strings.SplitSeq(httpResponse.Headers, "\r\n")
	for header := range splitHeaders {
		if strings.HasPrefix(header, title) {
			return true
		}
	}
	return false
}

func (httpResponse *HttpResponse) addHeader(title string, value string) {
	httpResponse.Headers += title + ": " + value + "\r\n"
}

func (httpResponse *HttpResponse) replaceHeader(title string, value string) {
	splitHeaders := strings.Split(httpResponse.Headers, "\r\n")
	for i, header := range splitHeaders {
		if strings.HasPrefix(header, title+":") {
			splitHeaders[i] = title + ": " + value
			httpResponse.Headers = strings.Join(splitHeaders, "\r\n")
			return
		}
	}
}

func (httpResponse *HttpResponse) addOrReplaceHeader(title string, value string) {
	if httpResponse.headerExists(title) {
		httpResponse.replaceHeader(title, value)
	} else {
		httpResponse.addHeader(title, value)
	}
}

func (httpResponse *HttpResponse) addOrReplaceHeaders(headers map[string]string) {
	for title, value := range headers {
		httpResponse.addOrReplaceHeader(title, value)
	}
}

func (httpResponse *HttpResponse) getHeader(title string) string {
	splitHeaders := strings.SplitSeq(httpResponse.Headers, "\r\n")
	for header := range splitHeaders {
		if after, ok := strings.CutPrefix(header, title+": "); ok {
			return after
		}
	}
	return ""
}

func (httpResponse *HttpResponse) compressBody(encoding string) ([]byte, error) {
	compressed, err := compressString(httpResponse.Body, encoding)
	return compressed, err
}

func (httpResponse *HttpResponse) handleFinalResponse(req *HttpRequest) {
	var contentType string
	if !httpResponse.headerExists("Content-Type") {
		contentType = "text/plain"
	} else {
		contentType = httpResponse.getHeader("Content-Type")
	}
	acceptEncodings := strings.SplitSeq(req.getHeader("Accept-Encoding"), ", ")
	for acceptEncoding := range acceptEncodings {
		if isEncodingSupported(acceptEncoding) {
			compressedBody, err := httpResponse.compressBody(acceptEncoding)
			if err != nil {
				fmt.Printf("Error encoding body in %s format.", acceptEncoding)
				continue
			}
			httpResponse.Body = string(compressedBody)
			httpResponse.addOrReplaceHeader("Content-Encoding", acceptEncoding)
			break
		}
	}
	httpResponse.addOrReplaceHeaders(map[string]string{
		"Content-Type":   contentType,
		"Content-Length": strconv.Itoa(len(httpResponse.Body)),
	})
	if req.getHeader("Connection") == "close" {
		httpResponse.addOrReplaceHeader("Connection", "close")
	}
}

func generic200Response() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       200,
		ReasonPhrase: "OK",
	}
}

func generic201Response() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       201,
		ReasonPhrase: "Created",
	}
}

func generic400Error() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       404,
		ReasonPhrase: "Not Found",
	}
}

func generic405Error() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       405,
		ReasonPhrase: "Method Not Allowed",
	}
}

func generic500Error() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       500,
		ReasonPhrase: "Internal Server Error",
	}
}
