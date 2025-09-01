package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Suppur/chirpy/internal/auth"
	"github.com/Suppur/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id         uuid.UUID     `json:"id"`
	Created_at time.Time     `json:"created_at"`
	Updated_at time.Time     `json:"updated_at"`
	Body       string        `json:"body"`
	User_id    uuid.NullUUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	tokin, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Please login before creating a chirp!", err)
		return
	}

	validID, err := auth.ValidateJWT(tokin, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error validating JWT", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	cleanedBody, err := badWordReplacement(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: uuid.NullUUID{
			UUID:  validID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID,
	})
}

func badWordReplacement(body string) (string, error) {
	lowerBody := strings.ToLower(body)
	hasBadWords := strings.Contains(lowerBody, "kerfuffle") ||
		strings.Contains(lowerBody, "sharbert") ||
		strings.Contains(lowerBody, "fornax")

	if !hasBadWords {
		return body, nil
	}

	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " "), nil
}
