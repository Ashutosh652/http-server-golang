package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func handleRoot(conn net.Conn, config Config, req *HttpRequest) *HttpResponse {
	return generic200Response()
}

func handleEcho(conn net.Conn, config Config, req *HttpRequest) *HttpResponse {
	splitString := strings.Split(req.Target, "/")
	pathString := splitString[len(splitString)-1]
	httpResponse := generic200Response()
	httpResponse.addOrReplaceHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(pathString)),
	})
	httpResponse.Body = pathString
	return httpResponse
}

func handleUserAgent(conn net.Conn, config Config, req *HttpRequest) *HttpResponse {
	userAgent := req.getHeader("User-Agent")
	httpResponse := generic200Response()
	httpResponse.addOrReplaceHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(userAgent)),
	})
	httpResponse.Body = userAgent
	return httpResponse
}

func handleFiles(conn net.Conn, config Config, req *HttpRequest) *HttpResponse {
	splitString := strings.Split(req.Target, "/")
	fileName := splitString[len(splitString)-1]
	httpResponse := generic200Response()

	switch req.Method {
	case HTTP_GET_METHOD:
		content, err := os.ReadFile(config.FileDirectory + fileName)
		if err != nil {
			return generic400Error()
		}
		httpResponse.addOrReplaceHeaders(map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": strconv.Itoa(len(content)),
		})
		httpResponse.Body = string(content)
		return httpResponse

	case HTTP_POST_METHOD:
		info, err := os.Stat(config.FileDirectory)
		if err != nil {
			httpResponse := generic400Error()
			if os.IsNotExist(err) {
				httpResponse.Body = "Directory does not exist."
			}
			return httpResponse
		}
		if !info.IsDir() {
			httpResponse := generic400Error()
			httpResponse.Body = fmt.Sprintf("%s is not a directory", config.FileDirectory)
			return httpResponse
		}

		fullPath := filepath.Join(config.FileDirectory, fileName)
		file, err := os.Create(fullPath)
		if err != nil {
			httpResponse := generic400Error()
			httpResponse.Body = "Could not create file."
			return httpResponse
		}
		defer file.Close()

		_, err = file.WriteString(req.Body)
		if err != nil {
			httpResponse := generic400Error()
			httpResponse.Body = "Could not write content to file."
			return httpResponse
		}

		return generic201Response()

	default:
		return generic500Error()
	}
}
