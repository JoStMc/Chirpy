package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	if err := cfg.dbQueries.ResetUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Erorr reseting users: %v", err))
		return
	} 
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database cleared\n"))
}
