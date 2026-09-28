package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"web_stok_barang/config"
	"web_stok_barang/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("rifsshikii1811")

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User // Menggunakan model User (atau Struct lokal)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	if user.Role == "" {
		user.Role = "petugas"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal memproses password", http.StatusInternalServerError)
		return
	}

	// PERBAIKAN: Menggunakan config.DBRAW (bukan DBRAWRAWRAW)
	query := "INSERT INTO users (username, password, role) VALUES (?, ?, ?)"
	_, err = config.DBRAW.Exec(query, user.Username, hashedPassword, user.Role)
	if err != nil {
		http.Error(w, "Gagal mendaftarkan user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Register user berhasil!"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	// PERBAIKAN: Menggunakan config.DBRAW
	query := "SELECT id, username, password, role FROM users WHERE username = ?"
	err = config.DBRAW.QueryRow(query, input.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		http.Error(w, "Username tidak ditemukan", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Password salah", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login berhasil!",
		"token":   tokenString,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, okID := r.Context().Value("user_id").(int)
	if !okID {
		// Fallback jika ID tersimpan sebagai float64 saat JWT di-parse
		if floatID, ok := r.Context().Value("user_id").(float64); ok {
			userID = int(floatID)
			okID = true
		}
	}

	if !okID {
		http.Error(w, "Unauthorized: Data profil tidak ditemukan", http.StatusUnauthorized)
		return
	}

	var user models.User
	// PERBAIKAN: Menggunakan config.DBRAW
	query := "SELECT id, username, role FROM users WHERE id = ?"
	err := config.DBRAW.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
