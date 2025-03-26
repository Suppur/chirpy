package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReqNumber(w http.ResponseWriter, r *http.Request) {
	hitCount := cfg.fileServerHits.Load()

	w.Header().Set("Content-Type:", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", hitCount)))

}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)

	w.Header().Set("Content-Type:", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
