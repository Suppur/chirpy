package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const metricz = `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`

func (cfg *apiConfig) handlerReqNumber(w http.ResponseWriter, r *http.Request) {
	hitCount := cfg.fileServerHits.Load()

	w.Header().Set("Content-Type:", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, metricz, hitCount)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)

	w.Header().Set("Content-Type:", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errP := returnVals{
			Error: fmt.Sprintf("Something went wrong %s", err),
		}
		log.Print(errP.Error)
		w.WriteHeader(400)
		return
	}

	if len(params.Body) > 140 {
		errP := returnVals{
			Error: "Chirp is too long",
		}
		log.Print(errP.Error)
		w.WriteHeader(400)
		return
	}

	respBody := returnVals{
		Valid: true,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling response body: %s", err)
		w.WriteHeader(400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
