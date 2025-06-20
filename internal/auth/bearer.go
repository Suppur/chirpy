package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		log.Print("error: token string is empty")
		return "", errors.New("token string cannot be empty")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	return tokenString, nil
}
