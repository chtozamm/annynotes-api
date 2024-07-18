package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("AUTH_SECRET_KEY"))
var tokens []string

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(id, email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: id,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT extracts "Authorization" header from the HTTP request, validates JWT and returns its claims.
func ValidateJWT(headers http.Header) (*Claims, error) {
	bearerToken := headers.Get("Authorization")
	authHeader := strings.Split(bearerToken, " ")
	if bearerToken == "" {
		return nil, errors.New("authorization token is required")
	}
	if len(authHeader) != 2 {
		return nil, fmt.Errorf(`invalid authorization token: expected format "Authorization: Bearer <token>, got: %v`, authHeader)
	}
	reqToken := authHeader[1]

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.New("token is expired")
		case errors.Is(err, jwt.ErrTokenInvalidClaims):
			return nil, errors.New("malformed token: token contains invalid claims")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("token is malformed")
		default:
			return nil, errors.New("couldn't parse the authorization token")
		}
	}
	if !tkn.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}
