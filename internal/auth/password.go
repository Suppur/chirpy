package auth

import (
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost, err := bcrypt.Cost([]byte(password))
	if err != nil {
		log.Printf("error retrieving hashing cost: %s", err)
	}
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Print(err)
	}
	return hex.EncodeToString(hashedPW[:]), nil
}

func CheckPasswordHash(hash, password string) error {
	decHash, err := hex.DecodeString(hash)
	if err != nil {
		log.Printf("error decoding hash: %s", err)
	}
	ok := bcrypt.CompareHashAndPassword(decHash, []byte(password))
	if ok != nil {
		log.Printf("error, incorrect password: %s", err)
	}
	return nil
}
