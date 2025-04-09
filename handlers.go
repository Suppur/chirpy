package main

import (
	"errors"
	"fmt"
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
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "endpoint not to be used in prod!", errors.New("forbidden"))
	}
	if err := cfg.db.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete users from db", err)
	}

	cfg.fileServerHits.Store(0)

	w.Header().Set("Content-Type:", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
