package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"web_stok_barang/config"
	"web_stok_barang/models"
)

// GetPengiriman - Mengambil seluruh data pengiriman untuk Dashboard & Tabel Petugas
func GetPengiriman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method tidak diizinkan",
		})
		return
	}

	var listPengiriman []models.Pengiriman

	err := config.DB.Preload("Barang").Order("created_at desc").Find(&listPengiriman).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Gagal mengambil data pengiriman",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listPengiriman)
}

// CreatePengiriman - Menambah data pengiriman baru
func CreatePengiriman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method tidak diizinkan",
		})
		return
	}

	var input models.Pengiriman
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Format input JSON tidak valid",
		})
		return
	}

	// VALIDASI: Cek apakah IDBarang sudah dipilih
	if input.IDBarang <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Barang harus dipilih terlebih dahulu (id_barang tidak boleh kosong)",
		})
		return
	}

	// Default status jika tidak diisi
	if input.Status == "" {
		input.Status = "Diproses"
	}

	if err := config.DB.Create(&input).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Gagal membuat data pengiriman",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pengiriman berhasil dibuat",
		"data":    input,
	})
}

// UpdatePengiriman - Mengubah status / konfirmasi pengiriman (/api/pengiriman/:id)
func UpdatePengiriman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method tidak diizinkan",
		})
		return
	}

	// Ambil ID Pengiriman dari URL Path
	pathParts := strings.Split(r.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ID Pengiriman tidak valid",
		})
		return
	}

	// 1. Cek dulu apakah data pengiriman ada dan barang sudah dipilih
	var pengiriman models.Pengiriman
	if err := config.DB.Where("id_pengiriman = ?", id).First(&pengiriman).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Data pengiriman tidak ditemukan",
		})
		return
	}

	// 2. Validasi apakah pengiriman tersebut memiliki barang yang valid
	if pengiriman.IDBarang <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal konfirmasi: Pengiriman ini belum memiliki barang yang dipilih",
		})
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Format data tidak valid",
		})
		return
	}

	// 3. PERBAIKAN: Gunakan kolom "id_pengiriman" bukan "id"
	if err := config.DB.Model(&models.Pengiriman{}).Where("id_pengiriman = ?", id).Update("status", req.Status).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Gagal memperbarui status pengiriman",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Status pengiriman berhasil diperbarui",
	})
}
