package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Expires_in int    `json:"expires_in_seconds,omitempty"`
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

	expiration := time.Duration(params.Expires_in) * time.Second
	if expiration == 0 {
		expiration = time.Hour
	}
	if expiration > time.Hour {
		expiration = time.Hour
	}

	token, err := auth.MakeJWT(usr.ID, cfg.secret, expiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create JWT", err)
		return
	}

	type UsrToken struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Token      string    `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, UsrToken{
		Id:         usr.ID,
		Created_at: usr.CreatedAt,
		Updated_at: usr.UpdatedAt,
		Email:      usr.Email,
		Token:      token,
	})

}
