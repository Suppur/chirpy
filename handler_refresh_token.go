package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
	"github.com/Suppur/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "invalid token", errors.New("token missing"))
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	validToken, err := cfg.db.GetToken(r.Context(), token)
	if err != nil || token != validToken.Token || validToken.ExpiresAt.Time.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	newToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error refreshing token", err)
		return
	}

	refreshedToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  newToken,
		UserID: validToken.UserID,
		ExpiresAt: sql.NullTime{
			Time: time.Now().Add(time.Hour * 1440),
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error inserting token in DB", err)
	}

	respondWithJSON(w, http.StatusOK, refreshedToken)
}
