package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JoStMc/Chirpy/internal/auth"
	"github.com/JoStMc/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type userRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
} 

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[userRequest](r)
	if err != nil {
	    respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %v", err))
		return
	} 

	if len(params.Email) < 1 {
	    respondWithError(w, http.StatusBadRequest, "Email too short")
		return
	} 

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %v", err))
		return
	} 

	hashedParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPass,
	} 
	res, err := cfg.dbQueries.CreateUser(context.Background(), hashedParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
		return
	} 

	response := User{
	    ID: res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Email: res.Email,
	} 

	respondWithJSON(w, http.StatusCreated, response)
} 

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[userRequest](r)
	if err != nil {
	    respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %v", err))
		return
	} 

	user, err := cfg.dbQueries.GetUser(context.Background(), params.Email)
	passwordMatches, _ := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !passwordMatches{
	    respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	} 

	response := User{
	    ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	} 
	respondWithJSON(w, http.StatusOK, response)
} 
