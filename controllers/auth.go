package controllers

import (
	"log"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"go-backend-technical-test/database"
	"go-backend-technical-test/models"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	queryExisting := "SELECT id, fullname, email FROM users WHERE email = ?"
	err := database.DB.QueryRow(queryExisting, user.Email).Scan(&user.ID, &user.Fullname, &user.Email)
	if err == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (fullname, email, password) VALUES (?, ?, ?)"
	_, err = database.DB.Exec(query, user.Fullname, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error saving user to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "User registered successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	query := "SELECT id, fullname, email, password FROM users WHERE email = ?"
	err := database.DB.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Fullname, &user.Email, &user.Password)
	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"fullname": user.Fullname,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenString,
	})
}
