package auth

import (
	"encoding/json"
	"net/http"
	
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
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
func Validate(w http.ResponseWriter, r *http.Request){
	tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
        http.Error(w, "Authorization header missing", http.StatusUnauthorized)
        return
    }
	_,err := ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Token is valid"))
}
