package handlers

import (
	"encoding/json"
	"fifthtask/internal/db"
	usertype "fifthtask/internal/user"
	"fifthtask/pkg/jwt"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var registerUser usertype.User
	err := json.NewDecoder(r.Body).Decode(&registerUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`INSERT INTO users(username, email) VALUES ($1, $2)`, registerUser.Username, registerUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var loginUser usertype.User
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow(`SELECT username, email FROM users WHERE username=$1 AND email=$2`, loginUser.Username, loginUser.Email)
	var username, email string
	err = row.Scan(&username, &email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	secretKey := []byte("your-secret-key")
	token, err := jwt.MadeToken(loginUser.Username, loginUser.Email, secretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	secretKey := []byte("your-secret-key")
	tokenString := r.Header.Get("Authorization")
	claims, err := jwt.CheckToken(tokenString, secretKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var profileUser usertype.User

	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	profileUser.Username = username

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	profileUser.Email = email

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profileUser)
}
