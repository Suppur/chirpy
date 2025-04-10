package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching chirps from database", err)
		return
	}

	returnChirps := []Chirp{}
	for _, chirp := range chirps {
		returnChirps = append(returnChirps, Chirp{
			Id:         chirp.ID,
			Created_at: chirp.CreatedAt,
			Updated_at: chirp.UpdatedAt,
			Body:       chirp.Body,
			User_id:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, returnChirps)
}

func (cfg *apiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	pathVal := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(pathVal)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid id, must be a valid UUID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error, chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		Body:       chirp.Body,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Id:         chirp.ID,
		User_id:    chirp.UserID,
	})
}
