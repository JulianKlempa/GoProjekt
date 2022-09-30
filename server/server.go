package server

import (
	"digitalDistribution/server/download"
	mainpage "digitalDistribution/server/mainPage"
	"digitalDistribution/server/upload"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func StartServer() {
	http.HandleFunc("/", mainpage.ServeHTTP)

	http.HandleFunc("/upload/", upload.ServeHTTP)

	http.HandleFunc("/downloads/", download.ServeHTTP)

	baseFilePath, err := os.Getwd()

	if err != nil {
		fmt.Println("Error while reading ssl-files")
		os.Exit(1)
	}

	fmt.Println(baseFilePath)

	certificateFilePath := filepath.Join(baseFilePath, "resources", "localhost.crt")
	keyFilePath := filepath.Join(baseFilePath, "resources", "localhost.key")

	log.Fatal(http.ListenAndServeTLS(":9000", certificateFilePath, keyFilePath, nil))
}
