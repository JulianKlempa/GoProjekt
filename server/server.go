package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello, World!")
	})
	log.Fatal(http.ListenAndServe(":9000", nil))
}
