//author Julian Klempa (4085242)

package authentication

import (
	"crypto/subtle"
	"digitalDistribution/configuration"
)

func Authenticate(username string, password [32]byte) bool {
	credentials := configuration.GetCredentials()
	pass := credentials[username]

	passwordMatch := (subtle.ConstantTimeCompare(password[:], pass[:]) == 1)
	return passwordMatch
}
