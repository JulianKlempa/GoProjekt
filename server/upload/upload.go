package upload

import (
	"crypto/sha256"
	"crypto/subtle"
	filemanager "digitalDistribution/fileManager"
	"fmt"
	"net/http"
)

type UploadHandler struct {
	Handler http.Handler
}

func (h UploadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()

	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte("admin"))
		expectedPasswordHash := sha256.Sum256([]byte("admin"))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
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
