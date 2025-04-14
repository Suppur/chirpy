package auth_test

import (
	"testing"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "suPerCooLStrongPW1!"

	hashed, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword return eerror: %v", err)
	}

	if hashed == "" {
		t.Fatal("HashPassword returned an empty string")
	}

	err = auth.CheckPasswordHash(hashed, password)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}

	wrongPassword := "notTheSamePassword&"
	err = auth.CheckPasswordHash(hashed, wrongPassword)
	if err == nil {
		t.Fatal("CheckPassword should have failed with the incorrect pw, but returned no error")
	}
}

func TestTokenCreateAndValidate(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret!!!"
	expiresIn := 60000000000
	expiredTest := -30000

	testToken, err := auth.MakeJWT(userID, tokenSecret, time.Duration(expiresIn))
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}

	expiredToken, err := auth.MakeJWT(userID, tokenSecret, time.Duration(expiredTest))
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}

	if testToken == "" {
		t.Fatal("token is empty")
	}

	validatedID, err := auth.ValidateJWT(testToken, tokenSecret)
	if err != nil {
		t.Fatalf("error validating token: %v", err)
	}

	if validatedID != userID {
		t.Fatal("mismatched user IDs")
	}

	expiredId, err := auth.ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Fatalf("error, expired token was validated: %v", expiredId)
	}
}
