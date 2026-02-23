package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)


func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
	respondWithJSON(w, http.StatusOK, struct{Valid bool `json:"valid"`}{Valid: true})
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
