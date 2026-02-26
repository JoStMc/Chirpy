package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/JoStMc/Chirpy/internal/auth"
)

type refreshResponse struct {
    Token string `json:"token"`
} 

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
	    respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	} 

	ctx := r.Context()
	userTokens, err := cfg.dbQueries.GetUserFromRefreshToken(ctx, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	} 

	if userTokens.ExpiresAt.Before(time.Now()) || userTokens.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token expired")
		return
	} 

	jwtoken, err := auth.MakeJWT(userTokens.UserID, cfg.jwtSecret, 3600 * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error generating token: %v", err))
		return
	} 
	response := refreshResponse{
	    Token: jwtoken,
	} 
	respondWithJSON(w, http.StatusOK, response)
} 


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
	    respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	} 

	ctx := r.Context()
	err = cfg.dbQueries.RevokeRefreshToken(ctx, token)
	if err != nil && err != sql.ErrNoRows {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error revoking token: %v", err))
		return
	} 

	w.WriteHeader(http.StatusNoContent)
} 
