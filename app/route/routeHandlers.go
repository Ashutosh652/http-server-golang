package route

import (
	"fmt"
	"http-server-go/app/config"
	"http-server-go/app/http"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func handleRoot(conn net.Conn, config config.Config, req *http.HttpRequest) *http.HttpResponse {
	return http.Generic200Response()
}

func handleEcho(conn net.Conn, config config.Config, req *http.HttpRequest) *http.HttpResponse {
	splitString := strings.Split(req.Target, "/")
	pathString := splitString[len(splitString)-1]
	httpResponse := http.Generic200Response()
	httpResponse.AddOrReplaceHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(pathString)),
	})
	httpResponse.Body = pathString
	return httpResponse
}

func handleUserAgent(conn net.Conn, config config.Config, req *http.HttpRequest) *http.HttpResponse {
	userAgent := req.GetHeader("User-Agent")
	httpResponse := http.Generic200Response()
	httpResponse.AddOrReplaceHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(userAgent)),
	})
	httpResponse.Body = userAgent
	return httpResponse
}

func handleFiles(conn net.Conn, config config.Config, req *http.HttpRequest) *http.HttpResponse {
	splitString := strings.Split(req.Target, "/")
	fileName := splitString[len(splitString)-1]
	httpResponse := http.Generic200Response()

	switch req.Method {
	case http.HTTP_GET_METHOD:
		content, err := os.ReadFile(config.FileDirectory + fileName)
		if err != nil {
			return http.Generic400Error()
		}
		httpResponse.AddOrReplaceHeaders(map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": strconv.Itoa(len(content)),
		})
		httpResponse.Body = string(content)
		return httpResponse

	case http.HTTP_POST_METHOD:
		info, err := os.Stat(config.FileDirectory)
		if err != nil {
			httpResponse := http.Generic400Error()
			if os.IsNotExist(err) {
				httpResponse.Body = "Directory does not exist."
			}
			return httpResponse
		}
		if !info.IsDir() {
			httpResponse := http.Generic400Error()
			httpResponse.Body = fmt.Sprintf("%s is not a directory", config.FileDirectory)
			return httpResponse
		}

		fullPath := filepath.Join(config.FileDirectory, fileName)
		file, err := os.Create(fullPath)
		if err != nil {
			httpResponse := http.Generic400Error()
			httpResponse.Body = "Could not create file."
			return httpResponse
		}
		defer file.Close()

		_, err = file.WriteString(req.Body)
		if err != nil {
			httpResponse := http.Generic400Error()
			httpResponse.Body = "Could not write content to file."
			return httpResponse
		}

		return http.Generic201Response()

	default:
		return http.Generic500Error()
	}
}
