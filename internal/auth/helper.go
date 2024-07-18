package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword encrypts password using the bcrypt hashing algorithm.
func HashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashedPasswordBytes), err
}

// checkPassword returns true if hashed and provided passwords match.
func CheckPassword(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func RespondWithJWT(w http.ResponseWriter, userID, userEmail string) {

	token, err := GenerateJWT(userID, userEmail)
	if err != nil {
		log.Printf("Failed to generate JWT: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(&struct {
		Token string `json:"token"`
	}{Token: token})
	if err != nil {
		log.Printf("Failed to marshal token: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
