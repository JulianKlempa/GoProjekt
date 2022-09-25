package server

import (
	"digitalDistribution/server/upload"
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
	fs := http.StripPrefix("/upload/", http.FileServer(http.Dir(".")))
	http.Handle("/upload/", upload.UploadHandler{Handler: fs})

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
