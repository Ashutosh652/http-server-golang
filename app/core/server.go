package core

import (
	"fmt"
	"http-server-go/app/config"
	"http-server-go/app/http"
	"http-server-go/app/route"
	"io"
	"net"
	"os"
	"slices"
	"strings"
)

func isMethodAllowed(method http.HttpMethod, route route.Route) bool {
	return slices.Contains(route.AllowedMethods, method)
}

func handleRequest(conn net.Conn, config config.Config, routes []route.Route) {
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client (EOF)")
				return
			}
			fmt.Println("Error reading request: ", err.Error())
			return
		}
		if n == 0 {
			fmt.Println("No data received")
			return
		}

		httpRequest := http.ParseRequest(buffer[:n])

		routeMatch := false
		for _, route := range routes {
			if (route.MatchExact && httpRequest.Target == route.Path) || (!route.MatchExact && strings.HasPrefix(httpRequest.Target, route.Path)) {
				if !isMethodAllowed(httpRequest.Method, route) {
					httpResponse := http.Generic405Error()
					httpResponse.AddOrReplaceHeader("Allow", strings.Join(route.AllowedMethods.AsStrings(), ", ")+"\r\n")
					conn.Write([]byte(httpResponse.CreateResponse()))
					return
				}
				httpResponse := route.Handler(conn, config, httpRequest)
				httpResponse.HandleFinalResponse(httpRequest)
				response := httpResponse.CreateResponse()
				conn.Write([]byte(response))
				routeMatch = true
				break
			}
		}
		if !routeMatch {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		if httpRequest.GetHeader("Connection") == "close" {
			conn.Close()
		}
	}
}

func StartServer(config config.Config, routes []route.Route, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn, config, routes)
	}
}
