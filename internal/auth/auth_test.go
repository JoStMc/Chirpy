package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	tokenSecret := "secretthing"
	userID := uuid.New()
	expirationTime := 2 * time.Second


	token, err := MakeJWT(userID, tokenSecret, expirationTime)
	if err != nil {
		t.Errorf("Error making JWT: %v", err)
	} 

	validCase, err := ValidateJWT(token, tokenSecret)
	if validCase != userID || err != nil {
		t.Errorf("Validate failed: %v", err)
	} 

	invalidCase, err := ValidateJWT(token, "differentThing")
	if invalidCase == userID || err == nil {
		t.Errorf("Validate failed: %v", err)
	} 

	time.Sleep(3 * time.Second)

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Errorf("Expected an error for an expired token")
	} 
}
