package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/JoStMc/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type upgradeUserRequest struct {
    Event string `json:"event"`
    Data struct {
		UserId uuid.UUID `json:"user_id"`
    } `json:"data"`
} 

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
	    respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	} 
	
	data, err := decodeJSON[upgradeUserRequest](r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	} 
	if data.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	} 

	err = cfg.dbQueries.UpgradeUserToChirpyRed(r.Context(), data.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		    respondWithError(w, http.StatusNotFound, "User not found")
			return
		} 
		respondWithError(w, http.StatusInternalServerError, "Could not upgrade user")
		return
	} 
	w.WriteHeader(http.StatusNoContent)
}
