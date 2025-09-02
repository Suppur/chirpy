package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
	"github.com/Suppur/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Usr
		Token         string `json:"token"`
		Refresh_token string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	usr, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found, please enter a valid email", err)
		return
	}

	err = auth.CheckPasswordHash(usr.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password, try again", err)
		return
	}

	token, err := auth.MakeJWT(usr.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create JWT", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create refresh token", err)
		return
	}

	db_token, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refresh_token,
		UserID: uuid.NullUUID{
			UUID:  usr.ID,
			Valid: true,
		},
		ExpiresAt: sql.NullTime{
			Time: time.Now().Add(time.Hour * 1440),
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating database entry for token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Usr: Usr{
			Id:         usr.ID,
			Created_at: usr.CreatedAt,
			Updated_at: usr.UpdatedAt,
			Email:      usr.Email,
		},
		Token:         token,
		Refresh_token: db_token,
	})

}
