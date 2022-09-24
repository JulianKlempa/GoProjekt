package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func StartServer() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello, World!")
	})
	baseFilePath, err := os.Getwd()

	if err != nil {
		fmt.Println("Error while reading ssl-files")
		os.Exit(1)
	}

	fmt.Println(baseFilePath)

	certificateFilePath := filepath.Join(baseFilePath, "server", "ssl", "localhost.crt")
	keyFilePath := filepath.Join(baseFilePath, "server", "ssl", "localhost.key")

	log.Fatal(http.ListenAndServeTLS(":9000", certificateFilePath, keyFilePath, nil))
}
