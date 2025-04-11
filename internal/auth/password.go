package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, err error) {
	hashedPW, err := bcrypt.GenerateFromPassword(password)
	if err != nil {
		log.Print(err)
	}
	return hashedPW, nil
}
