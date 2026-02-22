package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)


func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	} 

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %v\n", err))
		return
	} 

	if len(params.Body) > 140 {
		respondWithJSON(w, http.StatusBadRequest, struct{Err string `json:"error"`}{Err: "Chirp is too long"})
		return
	} 

	respondWithJSON(w, http.StatusOK, struct{Valid bool `json:"valid"`}{Valid: true})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Erorr: %s", msg)
	} 
	respondWithJSON(w, code, struct{Error string `json:"error"`}{Error: msg,} )
} 

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} 

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
} 
