package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
	"github.com/Suppur/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving bearer token", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
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
		return
	}

	respondWithJSON(w, http.StatusNoContent, err)
}
