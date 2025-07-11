package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"
)

func isMethodAllowed(method HttpMethod, route Route) bool {
	return slices.Contains(route.AllowedMethods, method)
}

func handleRequest(conn net.Conn, config Config, routes []Route) {
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

		httpRequest := parseRequest(buffer[:n])

		routeMatch := false
		for _, route := range routes {
			if (route.MatchExact && httpRequest.Target == route.Path) || (!route.MatchExact && strings.HasPrefix(httpRequest.Target, route.Path)) {
				if !isMethodAllowed(httpRequest.Method, route) {
					httpResponse := generic405Error()
					httpResponse.addOrReplaceHeader("Allow", strings.Join(route.AllowedMethods.AsStrings(), ", ")+"\r\n")
					fmt.Println(httpResponse.createResponse())
					conn.Write([]byte(httpResponse.createResponse()))
					return
				}
				httpResponse := route.Handler(conn, config, httpRequest)
				httpResponse.handleFinalResponse(httpRequest)
				response := httpResponse.createResponse()
				fmt.Println(response)
				conn.Write([]byte(response))
				routeMatch = true
				break
			}
		}
		if !routeMatch {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		if httpRequest.getHeader("Connection") == "close" {
			conn.Close()
		}
	}
}

func StartServer(config Config, routes []Route) {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
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
