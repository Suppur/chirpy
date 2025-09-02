package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Println("error signing token")
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Println("error validating token")
		return uuid.Nil, err
	}
	subject, err := parsedToken.Claims.GetSubject()
	if err != nil {
		log.Println("error retrieving token claims")
		return uuid.Nil, err
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		log.Println("error retrieving userID")
		return uuid.Nil, err
	}

	return userID, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)

	rand.Read(key)

	encoded_str := hex.EncodeToString(key)
	if encoded_str == "" {
		log.Println("empty encoded string")
		return "", errors.New("error generating refresh token")
	}

	return encoded_str, nil

}
