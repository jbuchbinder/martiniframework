package martiniframework

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

// Token is the authorization token extracted from the request.
type Token string

// TokenFunc returns a Handler that authenticates via an Authentication: Bearer header using the provided function.
// The function should return true for a valid token.
func TokenFunc(authfn func(string) (bool, SessionModel)) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		log.Print("TokenFunc()")
		auth := req.Header.Get("Authorization")
		log.Print("TokenFunc(): Authorization: " + auth)
		if len(auth) < 7 || auth[:7] != "Bearer " {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return
		}
		authed, ses := authfn(auth[7:])
		if !authed {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return
		}

		// Inject Token
		c.Map(Token(auth[7:]))

		// Inject SessionModel
		c.Map(ses)
	}
}
