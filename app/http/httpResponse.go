package http

import (
	"fmt"
	"http-server-go/app/utils"
	"strconv"
	"strings"
)

type HttpResponse struct {
	Version      string
	Status       int
	ReasonPhrase string
	BaseHttp
}

func (httpResponse *HttpResponse) CreateResponse() string {
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s", httpResponse.Version, httpResponse.Status, httpResponse.ReasonPhrase, httpResponse.Headers, httpResponse.Body)
}

func (http *BaseHttp) headerExists(title string) bool {
	splitHeaders := strings.SplitSeq(http.Headers, "\r\n")
	for header := range splitHeaders {
		if strings.HasPrefix(header, title) {
			return true
		}
	}
	return false
}

func (http *BaseHttp) addHeader(title string, value string) {
	http.Headers += title + ": " + value + "\r\n"
}

func (http *BaseHttp) replaceHeader(title string, value string) {
	splitHeaders := strings.Split(http.Headers, "\r\n")
	for i, header := range splitHeaders {
		if strings.HasPrefix(header, title+":") {
			splitHeaders[i] = title + ": " + value
			http.Headers = strings.Join(splitHeaders, "\r\n")
			return
		}
	}
}

func (http *BaseHttp) AddOrReplaceHeader(title string, value string) {
	if http.headerExists(title) {
		http.replaceHeader(title, value)
	} else {
		http.addHeader(title, value)
	}
}

func (http *BaseHttp) AddOrReplaceHeaders(headers map[string]string) {
	for title, value := range headers {
		http.AddOrReplaceHeader(title, value)
	}
}

func (httpResponse *HttpResponse) CompressBody(encoding string) ([]byte, error) {
	compressed, err := utils.CompressString(httpResponse.Body, encoding)
	return compressed, err
}

func (httpResponse *HttpResponse) HandleFinalResponse(req *HttpRequest) {
	var contentType string
	if !httpResponse.headerExists("Content-Type") {
		contentType = "text/plain"
	} else {
		contentType = httpResponse.GetHeader("Content-Type")
	}
	acceptEncodings := strings.SplitSeq(req.GetHeader("Accept-Encoding"), ", ")
	for acceptEncoding := range acceptEncodings {
		if utils.IsEncodingSupported(acceptEncoding) {
			compressedBody, err := httpResponse.CompressBody(acceptEncoding)
			if err != nil {
				fmt.Printf("Error encoding body in %s format.", acceptEncoding)
				continue
			}
			httpResponse.Body = string(compressedBody)
			httpResponse.AddOrReplaceHeader("Content-Encoding", acceptEncoding)
			break
		}
	}
	httpResponse.AddOrReplaceHeaders(map[string]string{
		"Content-Type":   contentType,
		"Content-Length": strconv.Itoa(len(httpResponse.Body)),
	})
	if req.GetHeader("Connection") == "close" {
		httpResponse.AddOrReplaceHeader("Connection", "close")
	}
}

func Generic200Response() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       200,
		ReasonPhrase: "OK",
	}
}

func Generic201Response() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       201,
		ReasonPhrase: "Created",
	}
}

func Generic400Error() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       404,
		ReasonPhrase: "Not Found",
	}
}

func Generic405Error() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       405,
		ReasonPhrase: "Method Not Allowed",
	}
}

func Generic500Error() *HttpResponse {
	return &HttpResponse{
		Version:      "HTTP/1.1",
		Status:       500,
		ReasonPhrase: "Internal Server Error",
	}
}
