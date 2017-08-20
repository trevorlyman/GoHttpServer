package main

import (
	"./HttpServer"
)

func main() {
	server := HttpServer.HttpServer{}
	server.Run("8080")
}
