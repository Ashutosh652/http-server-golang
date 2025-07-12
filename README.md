This is a simple custom HTTP server built in Go with the help of [CodeCrafters](https://app.codecrafters.io/courses/http-server/overview).

This project implements a basic HTTP server in Go without using high-level HTTP libraries. It demonstrates fundamental networking concepts and HTTP protocol implementation by building the server from the ground up using Go's standard library.


## Features

Custom HTTP Implementation: Built without using Go's net/http package
Low-level Socket Programming: Direct TCP socket handling
HTTP/1.1 Protocol Support: Implements core HTTP/1.1 features
Request Parsing: Manual parsing of HTTP requests
Response Generation: Custom HTTP response formatting
Concurrent Connections: Handles multiple client connections


## To run locally:
- Clone the repository.
- `cd http-server-golang`.
- `go run main.go`


## Test the server using curl:
- `curl http://localhost:4221/`
- `curl http://localhost:4221/echo/hello`
- `curl http://localhost:4221/user-agent`
