package upload

import (
	"crypto/sha256"
	filemanager "digitalDistribution/fileManager"
	"digitalDistribution/server/authentication"
	"fmt"
	"net/http"
)

type UploadHandler struct {
	Handler http.Handler
}

func (h UploadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()

	if ok {
		passwordHash := sha256.Sum256([]byte(password))

		if authentication.Authenticate(username, passwordHash) {
			fmt.Println("Authenticated")
			if req.Method != "POST" {
				fmt.Println("No POST")
				return
			}
			fmt.Println("Authenticated and post")

			req.ParseMultipartForm(32 << 20)
			file, header, err := req.FormFile("fileupload")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			filemanager.SaveFile(file, header.Filename)

			return

		}
	}
	res.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(res, "Unauthorized", http.StatusUnauthorized)
}
