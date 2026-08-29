package controllers

import (
	"encoding/json"
	"net/http"
	"web_stok_barang/config"
	"web_stok_barang/models"
)

func GetBarang(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id_barang, kategori, stok, harga, nomor_resi FROM barang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var listBarang []models.Barang
	for rows.Next() {
		var b models.Barang
		rows.Scan(&b.IDBarang, &b.Kategori, &b.Stok, &b.Harga, &b.NomorResi)
		listBarang = append(listBarang, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listBarang)
}

func CreateBarang(w http.ResponseWriter, r *http.Request) {
	var b models.Barang
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO barang (kategori, stok, harga, nomor_resi) VALUES (?, ?, ?, ?)"
	_, err = config.DB.Exec(query, b.Kategori, b.Stok, b.Harga, b.NomorResi)
	if err != nil {
		http.Error(w, "Gagal menambah barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Barang berhasil ditambahkan!"})
}

func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID barang diperlukan", http.StatusBadRequest)
		return
	}

	var b models.Barang
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	query := "UPDATE barang SET kategori = ?, stok = ?, harga = ?, nomor_resi = ? WHERE id_barang = ?"
	_, err = config.DB.Exec(query, b.Kategori, b.Stok, b.Harga, b.NomorResi, id)
	if err != nil {
		http.Error(w, "Gagal mengupdate barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Barang berhasil diupdate!"})
}

func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID barang diperlukan", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM barang WHERE id_barang = ?"
	_, err := config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Gagal menghapus barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Barang berhasil dihapus!"})
}
