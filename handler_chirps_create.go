package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Suppur/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	type returnVals struct {
		Id         uuid.UUID     `json:"id"`
		Created_at time.Time     `json:"created_at"`
		Updated_at time.Time     `json:"updated_at"`
		Body       string        `json:"body"`
		User_id    uuid.NullUUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
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
			UUID:  params.User_id,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp", err)
	}

	respBody := returnVals{
		Id:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, respBody)
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
