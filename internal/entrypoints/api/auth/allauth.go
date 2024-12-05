package auth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"nbox/internal/application"
	"nbox/internal/entrypoints/api/response"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationMiddleware interface{}

type Authn struct {
	credentials map[string]string
	config      *application.Config
}

func NewAuthn(prefix string, config *application.Config) *Authn {
	credentials := map[string]string{}
	_ = json.Unmarshal([]byte(os.Getenv(prefix)), &credentials)
	return &Authn{credentials: credentials, config: config}
}

func (a Authn) tryJwt(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := r.Context()

	//authHeader := r.Header.Get("Authorization")
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.config.HmacSecretKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return ctx, errors.New("invalid token")
	}

	expirationTime, _ := claims.GetExpirationTime()

	if expirationTime.Unix() < time.Now().Unix() {
		return ctx, errors.New("expired token")
	}

	ctx = context.WithValue(ctx, application.RequestUserName, claims["username"].(string))
	return ctx, nil
}

func (a Authn) tryBasicAuth(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := r.Context()
	realm := "api"
	user, pass, ok := r.BasicAuth()
	if !ok {
		unauthorized(w, realm)
		return ctx, errors.New("unauthorized")
	}

	credPass, credUserOk := a.credentials[user]
	if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
		unauthorized(w, realm)
		return ctx, errors.New("unauthorized")
	}

	ctx = context.WithValue(ctx, application.RequestUserName, user)
	return ctx, nil
}

func (a Authn) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//ctx := r.Context()
			authorization := r.Header.Get("Authorization")
			isBasic := strings.HasPrefix(strings.ToLower(authorization), "basic")
			isBearer := strings.HasPrefix(strings.ToLower(authorization), "bearer")

			if !isBearer && !isBasic {
				response.Error(w, r, errors.New("invalid auth method"), http.StatusUnauthorized)
				return
			}

			if isBasic {
				newCtx, err := a.tryBasicAuth(w, r)
				if err != nil {
					response.Error(w, r, err, http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r.WithContext(newCtx))
				return
			}

			if isBearer {
				newCtx, err := a.tryJwt(w, r)
				if err != nil {
					response.Error(w, r, err, http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r.WithContext(newCtx))
				return
			}

			//log.Printf("isBasic=%v | isBearer=%v \n", isBasic, isBearer)
			//next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
