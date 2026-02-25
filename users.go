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

type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type createUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
} 

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[createUserRequest](r)
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

	response := createUserResponse{
	    ID: res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Email: res.Email,
	} 

	respondWithJSON(w, http.StatusCreated, response)
} 

type loginResponse struct {
	ID        	 uuid.UUID `json:"id"`
	CreatedAt 	 time.Time `json:"created_at"`
	UpdatedAt 	 time.Time `json:"updated_at"`
	Email     	 string    `json:"email"`
	Token     	 string    `json:"token"`
	RefreshToken string	   `json:"refresh_token"`
}

type loginRequest struct {
	Email 			 string `json:"email"`
	Password 		 string `json:"password"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[loginRequest](r)
	if err != nil {
	    respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %v", err))
		return
	} 
	ctx := context.Background()

	user, err := cfg.dbQueries.GetUser(ctx, params.Email)
	passwordMatches, _ := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !passwordMatches{
	    respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	} 

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, 3600 * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to create token: %v", err))
		return
	} 

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.dbQueries.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating refresh token: %v", err))
		return
	} 

	response := loginResponse{
	    ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken,
	} 
	respondWithJSON(w, http.StatusOK, response)
} 
