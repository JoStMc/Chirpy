package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
} 

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
} 


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		Subject: userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims) 
	return token.SignedString([]byte(tokenSecret))
} 

type CustomClaims struct {
	jwt.RegisteredClaims
} 

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
	    return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userIDString, err := claims.GetSubject()
	if err != nil {
	    return uuid.Nil, err
	} 
	return uuid.Parse(userIDString)
} 


func GetBearerToken(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	if authString == "" {
	    return "", fmt.Errorf("authorization token not found")
	} 

	splitAuthString := strings.Fields(authString)
	if splitAuthString[0] != "Bearer" || len(splitAuthString) != 2 {
	    return "", fmt.Errorf("authorization token not found")
	} 

	return splitAuthString[1], nil
}
