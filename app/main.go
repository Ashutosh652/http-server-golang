package main

import (
	"http-server-go/app/config"
	"http-server-go/app/core"
	"http-server-go/app/route"
)

func main() {
	configuration := config.ParseConfig()
	routes := route.InitRoutes()
	core.StartServer(configuration, routes, 4221)
}
