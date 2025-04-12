package main

import (
	"encoding/json"
	"net/http"

	"github.com/Suppur/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	}

	respondWithJSON(w, http.StatusOK, Usr{
		Id:         usr.ID,
		Created_at: usr.CreatedAt,
		Updated_at: usr.UpdatedAt,
		Email:      usr.Email,
	})

}
