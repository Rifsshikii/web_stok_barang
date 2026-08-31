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
	var admin models.Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	if admin.Role == "" {
		admin.Role = "admin"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal memproses password", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO admin (nama, nomor_hp, username, password, role) VALUES (?, ?, ?, ?, ?)"
	_, err = config.DB.Exec(query, admin.Nama, admin.NomorHP, admin.Username, hashedPassword, admin.Role)
	if err != nil {
		http.Error(w, "Gagal mendaftarkan admin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Register admin berhasil!"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.Admin
	var admin models.Admin

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	query := "SELECT id_admin, nama, username, password, role FROM admin WHERE username = ?"
	err = config.DB.QueryRow(query, input.Username).Scan(&admin.IDAdmin, &admin.Nama, &admin.Username, &admin.Password, &admin.Role)
	if err != nil {
		http.Error(w, "Username tidak ditemukan", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Password salah", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"id_admin": admin.IDAdmin,
		"nama":     admin.Nama,
		"role":     admin.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login berhasil!",
		"token":   tokenString,
		"admin": map[string]interface{}{
			"id_admin": admin.IDAdmin,
			"nama":     admin.Nama,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
}
