package auth_test

import (
	"testing"

	"github.com/Suppur/chirpy/internal/auth"
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
