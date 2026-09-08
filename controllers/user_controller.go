package controllers

import (
	"encoding/json"
	"net/http"
	"web_stok_barang/config"
)

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query("SELECT id, username, role FROM users ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Gagal mengambil data user", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var u UserResponse
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
			http.Error(w, "Gagal membaca data user", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}
