package upload

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
	"net/http"
	"strings"
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

			req.ParseMultipartForm(32 << 20) // limit your max input length!
			var buf bytes.Buffer
			// in your case file would be fileupload
			file, header, err := req.FormFile("file")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			name := strings.Split(header.Filename, ".")
			fmt.Printf("File name %s\n", name[0])
			// Copy the file data to my buffer
			io.Copy(&buf, file)
			// do something with the contents...
			// I normally have a struct defined and unmarshal into a struct, but this will
			// work as an example
			contents := buf.String()
			fmt.Println(contents)
			// I reset the buffer in case I want to use it again
			// reduces memory allocations in more intense projects
			buf.Reset()
			// do something else
			// etc write header
			return
		}
	}
	res.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(res, "Unauthorized", http.StatusUnauthorized)
}
