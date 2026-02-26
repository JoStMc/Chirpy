package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JoStMc/Chirpy/internal/auth"
	"github.com/JoStMc/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
} 

type createChirpRequests struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[createChirpRequests](r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %v", err))
		return
	} 

	if len(params.Body) > 140 {
		respondWithJSON(w, http.StatusBadRequest, struct{Err string `json:"error"`}{Err: "Chirp is too long"})
		return
	} 
	if len(params.Body) < 1 {
	    respondWithError(w, http.StatusBadRequest,	"Chirp must be at least 1 character")
		return
	} 
	params.Body = replaceBadWords(params.Body)

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
	    respondWithError(w, http.StatusUnauthorized, fmt.Sprint(err))
		return
	} 
	bearerID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Access denied: %v", err))
	} 
	if bearerID != params.UserID {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Access denied"))
		return
	} 

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams(params))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %v", err))
		return
	} 

	respondWithJSON(w, http.StatusCreated, Chirp(chirp))
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving chirps: %v", err))
		return
	} 

	apiChirps := make([]Chirp, len(chirps))
	for i, c := range chirps {
	    apiChirps[i] = Chirp(c)
	} 
	respondWithJSON(w, http.StatusOK, apiChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	} 

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving chirps: %v", err))
		return
	} 

	apiChirp := Chirp(chirp)
	respondWithJSON(w, http.StatusOK, apiChirp)
}

var badWords = map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
} 

func replaceBadWords(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		if _, ok := badWords[word]; ok {
			words[i] = "****"
		} 
	} 
	return strings.Join(words, " ")
} 


func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
	    respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	} 
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpId)
	if err != nil {
	    respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	} 
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting user token: %v", err))
		return
	} 
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
	    respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	} 

	if chirp.UserID != userID {
	    respondWithError(w, http.StatusForbidden, "cannot delete chirp")
		return
	} 

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting chirp: %v", err))
		return
	} 

	w.WriteHeader(http.StatusNoContent)
} 
