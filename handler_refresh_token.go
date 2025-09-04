package main

import (
	"net/http"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error retrieving bearer token", err)
		return
	}
	validToken, err := cfg.db.GetToken(r.Context(), token)
	if err != nil || validToken.ExpiresAt.Time.Before(time.Now()) || validToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	newToken, err := auth.MakeJWT(validToken.UserID.UUID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating new access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}
