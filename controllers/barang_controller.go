package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"web_stok_barang/config"
	"web_stok_barang/models"
)

func GetBarang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := config.DB.Query("SELECT id_barang, kode_barang, kategori, stok, harga, nomor_resi, status FROM barang")
	if err != nil {
		http.Error(w, "Gagal mengambil data barang: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var listBarang []models.Barang
	for rows.Next() {
		var b models.Barang
		if err := rows.Scan(&b.IDBarang, &b.KodeBarang, &b.Kategori, &b.Stok, &b.Harga, &b.NomorResi, &b.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		listBarang = append(listBarang, b)
	}

	json.NewEncoder(w).Encode(listBarang)
}

func CreateBarang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b models.Barang

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO barang (kode_barang, kategori, stok, harga, nomor_resi, status) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := config.DB.Exec(query, b.KodeBarang, b.Kategori, b.Stok, b.Harga, b.NomorResi, b.Status)
	if err != nil {
		http.Error(w, "Gagal menambah barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	b.IDBarang = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Barang berhasil ditambahkan!",
		"data":    b,
	})
}

func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/barang/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID barang tidak valid", http.StatusBadRequest)
		return
	}

	var b models.Barang
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	query := "UPDATE barang SET kode_barang=?, kategori=?, stok=?, harga=?, nomor_resi=?, status=? WHERE id_barang=?"
	_, err = config.DB.Exec(query, b.KodeBarang, b.Kategori, b.Stok, b.Harga, b.NomorResi, b.Status, id)
	if err != nil {
		http.Error(w, "Gagal mengupdate barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	b.IDBarang = id
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Barang berhasil diperbarui!",
		"data":    b,
	})
}

func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/barang/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID barang tidak valid", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM barang WHERE id_barang=?"
	_, err = config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Gagal menghapus barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Barang berhasil dihapus!",
	})
}
