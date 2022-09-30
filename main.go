package main

import (
	"crypto/sha256"
	"digitalDistribution/configuration"
	filemanager "digitalDistribution/fileManager"
	server "digitalDistribution/server"
	"flag"
)

func main() {
	var username string
	flag.StringVar(&username, "username", "admin", "username for https server")
	var password string
	flag.StringVar(&password, "password", "admin", "password for https server")
	passwordHash := sha256.Sum256([]byte(password))
	var savesCount int
	flag.IntVar(&savesCount, "savesCount", 2, "naumber of saved states")
	var port string
	flag.StringVar(&port, "port", "9000", "the port for the server")

	credentials := make(map[string][]byte)
	credentials[username] = passwordHash[:]

	configuration.SetConfig(credentials, savesCount, port)
	filemanager.Setup()
	server.StartServer()
}
