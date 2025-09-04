package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving bearer token", err)
		return
	}

	validToken, err := cfg.db.GetToken(r.Context(), token)
	fmt.Println("Expires at:", validToken.ExpiresAt)
	fmt.Println("REvoked at:", validToken.RevokedAt)
	if err != nil || token != validToken.Token || validToken.ExpiresAt.Time.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	newToken, err := auth.MakeJWT(validToken.UserID.UUID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating new access token", err)
	}

	respondWithJSON(w, http.StatusOK, newToken)
}
