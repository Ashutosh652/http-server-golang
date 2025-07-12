A simple custom HTTP server built from scratch in Go, developed as part of the [CodeCrafters](https://app.codecrafters.io/courses/http-server/overview) "Build Your Own HTTP Server" challenge.

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
    - To enable file serving from a specific directory, use the --directory flag:
        `go run main.go --directory /tmp/`


## Test the server using curl:
- `curl http://localhost:4221/`
- `curl http://localhost:4221/echo/hello`
- `curl http://localhost:4221/user-agent`
- To test the file retrieving endpoint, you have to start the server with the `--directory` flag as mentioned above. Then create a test file in that same directory and provide the filename as a path parameter.
    ```bash
  echo "Hello, World!" > /tmp/test.txt
  curl http://localhost:4221/files/test.txt
  curl http://localhost:4221/files/nonexistent.txt
    ```
    The /files endpoint will:
  
      - Return the file contents as bytes if the file exists.
      - Return a 404 status code if the file doesn't exist.
      - Only serve files from the directory specified with --directory flag.
