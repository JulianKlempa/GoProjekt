package main

import (
	filemanager "digitalDistribution/fileManager"
	server "digitalDistribution/server"
)

func main() {
	filemanager.Setup()
	server.StartServer()
}
