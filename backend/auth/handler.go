package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type Credentials struct {
	Email    string `json:"email"`
	
}


// Login handles user login, generates JWT upon successful authentication.
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}


	token, err := GenerateJWT(creds.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	})

	w.Write([]byte("Login successful!"))
}
