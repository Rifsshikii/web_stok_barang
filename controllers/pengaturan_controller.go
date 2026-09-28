package controllers

import (
	"encoding/json"
	"net/http"

	"web_stok_barang/config"
)

type PengaturanDTO struct {
	ID           int    `json:"id"`
	NamaAplikasi string `json:"nama_aplikasi"`
	Alamat       string `json:"alamat"`
	Telepon      string `json:"telepon"`
}

// GET: Ambil data pengaturan
func GetPengaturan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT id, nama_aplikasi, COALESCE(alamat, ''), COALESCE(telepon, '') FROM pengaturan LIMIT 1`

	var p PengaturanDTO
	err := config.DBRAW.QueryRow(query).Scan(&p.ID, &p.NamaAplikasi, &p.Alamat, &p.Telepon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(p)
}

// PUT / POST: Update data pengaturan
func UpdatePengaturan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p PengaturanDTO
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE pengaturan SET nama_aplikasi = ?, alamat = ?, telepon = ? WHERE id = ?`
	_, err = config.DBRAW.Exec(query, p.NamaAplikasi, p.Alamat, p.Telepon, p.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Pengaturan berhasil diperbarui"})
}
