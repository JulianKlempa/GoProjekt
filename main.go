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
	flag.IntVar(&savesCount, "savesCount", 20, "naumber of saved states")
	if savesCount > 20 {
		savesCount = 20
	} else if savesCount < 0 {
		savesCount = 0
	}

	credentials := make(map[string][]byte)
	credentials[username] = passwordHash[:]

	conf := configuration.Configuration{}
	conf.Credentials = credentials
	conf.SavesCount = savesCount

	configuration.WriteConfig(conf)
	filemanager.Setup()
	server.StartServer()
}
