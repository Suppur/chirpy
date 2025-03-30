package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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

	respBody := returnVals{
		CleanedBody: cleanedBody,
	}

	respondWithJSON(w, http.StatusOK, respBody)
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
