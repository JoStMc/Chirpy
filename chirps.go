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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Error creating chirp: %v", err))
	} 


	respondWithJSON(w, http.StatusCreated, Chirp(chirp))
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
