package auth_test

import (
	"testing"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
	"github.com/google/uuid"
)

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
