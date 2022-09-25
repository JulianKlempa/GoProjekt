package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
)

type myhandler struct {
	fs http.Handler
}

func (h myhandler) ServeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		return
	}
	username, password, ok := r.BasicAuth()

	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte("admin"))
		expectedPasswordHash := sha256.Sum256([]byte("admin"))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
			fmt.Println("Authenticated")
			return
		}
	}
}
