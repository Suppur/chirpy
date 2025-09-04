package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Suppur/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "invalid token", errors.New("token missing"))
		return
	}

	err := cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token: token,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error, failed to revoke token", err)
	}

	respondWithJSON(w, http.StatusNoContent, err)
}
