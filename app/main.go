package main

func main() {
	configuration := ParseConfig()
	routes := initRoutes()
	StartServer(configuration, routes)
}
