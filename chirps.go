package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JoStMc/Chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	} 

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %v", err))
		return
	} 

	if len(params.Body) > 140 {
		respondWithJSON(w, http.StatusBadRequest, struct{Err string `json:"error"`}{Err: "Chirp is too long"})
		return
	} 
	params.Body = replaceBadWords(params.Body)

	chirp, err := cfg.dbQueries.CreateChirp(context.Background(), database.CreateChirpParams(params))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %v", err))
		return
	} 


	respondWithJSON(w, http.StatusCreated, Chirp(chirp))
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(context.Background())
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	} 

	chirp, err := cfg.dbQueries.GetChirps(context.Background(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving chirps: %v", err))
		return
	} 

	apiChirp := Chirp(chirp)
	respondWithJSON(w, http.StatusOK, apiChirp)
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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
