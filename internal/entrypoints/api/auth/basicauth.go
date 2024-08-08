package auth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"nbox/internal/application"
	"net/http"
	"os"
)

// basicAuth implements a simple middleware handler for adding basic http auth to a route.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
func basicAuth(realm string, credentials map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, pass, ok := r.BasicAuth()
			if !ok {
				unauthorized(w, realm)
				return
			}

			credPass, credUserOk := credentials[user]
			if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
				unauthorized(w, realm)
				return
			}

			ctx = context.WithValue(ctx, application.RequestUserName, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorized(w http.ResponseWriter, realm string) {
	log.Printf("unauthorized: realm: %s", realm)
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}

// NewBasicAuthFromEnv reads a set of credentials in from environment variables in
// the format {"user":"pass"} and returns
// middleware that will validate incoming requests.
func NewBasicAuthFromEnv(realm, prefix string) func(http.Handler) http.Handler {
	credentials := map[string]string{}

	err := json.Unmarshal([]byte(os.Getenv(prefix)), &credentials)
	if err != nil {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("Err Couldn't unmarshal %s. %v\n", prefix, err)
				unauthorized(w, realm)
			})
		}
	}
	return basicAuth(realm, credentials)
}
