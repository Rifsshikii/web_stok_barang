package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"web_stok_barang/config"
)

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// GetUsers - Ambil seluruh data user (Khusus Admin)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DBRAW.Query("SELECT id, username, role FROM users ORDER BY id DESC")
	if err != nil {
		http.Error(w, `{"error": "Gagal mengambil data user"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var u UserResponse
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
			http.Error(w, `{"error": "Gagal membaca data user"}`, http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}

// DeleteUser - Hapus user berdasarkan ID (/api/users/id)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Ambil ID dari URL Path (misal: /api/users/5 -> id = "5")
	pathParts := strings.Split(r.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID User tidak valid"})
		return
	}

	// 2. Eksekusi query DELETE ke database MySQL
	result, err := config.DBRAW.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal menghapus user. Pastikan data pesanan/transaksi user ini sudah tidak ada (Foreign Key constraint).",
		})
		return
	}

	// 3. Cek apakah ada baris data yang terpengaruh/terhapus
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User tidak ditemukan"})
		return
	}

	// 4. Berikan respon sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User berhasil dihapus",
	})
}
