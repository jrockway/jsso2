package main

import "github.com/jrockway/opinionated-server/server"

func main() {
	server.AppName = "jsso2"
	server.Setup()
	server.ListenAndServe()
}
