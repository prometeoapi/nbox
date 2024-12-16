package auth

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"nbox/internal/application"
	"nbox/internal/entrypoints/api/response"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//var secretKey = []byte("5e78e2ff5f63245c2ebd078b89eaf84085e378b0d24bcc927291b0fcc66baffc")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenHandler JWT
// @Summary Token
// @Description authentication token
// @Tags auth
// @Accept       json
// @Produce      json
// @Param data body TokenRequest true "Payload"
// @Success 200 {object} object{token=string} "Token generated successfully"
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/auth/token [post]
func (a Authn) TokenHandler(w http.ResponseWriter, r *http.Request) {
	credentials := map[string]string{}
	err := json.Unmarshal([]byte(os.Getenv(application.PrefixBasicAuthCredentials)), &credentials)
	if err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	payload := &TokenRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	credPass, credUserOk := credentials[payload.Username]
	if !credUserOk || subtle.ConstantTimeCompare([]byte(payload.Password), []byte(credPass)) != 1 {
		response.Error(w, r, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	now := time.Now()
	claims := &Claims{
		Username: payload.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Crear el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"username": payload.Username,
	//	"exp":      time.Now().Add(time.Hour * 24).Unix(),
	//	"iat":      time.Now().Unix(),
	//	"nbf":      time.Now().Unix(),
	//})

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(a.config.HmacSecretKey)
	if err != nil {
		//http.Error(w, "Error generating token", http.StatusInternalServerError)
		response.Error(w, r, err, http.StatusInternalServerError)
		return
	}

	response.Success(w, r, map[string]string{"token": tokenString})
}
