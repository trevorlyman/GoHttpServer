package main

import (
	"./HttpServer"
)

func main() {
	port := "8080"
	server := HttpServer.HttpServer{}
	server.Run(port)
}
