package auth

import (
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return hex.EncodeToString(hashedPW[:]), nil
}

func CheckPasswordHash(hash, password string) error {
	decHash, err := hex.DecodeString(hash)
	if err != nil {
		log.Printf("error decoding hash: %v", err)
		return err
	}
	err = bcrypt.CompareHashAndPassword(decHash, []byte(password))
	if err != nil {
		log.Printf("error, incorrect password: %v", err)
		return err
	}
	return nil
}
