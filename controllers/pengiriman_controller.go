package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"web_stok_barang/config"
	"web_stok_barang/models"
)

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

	if input.IDBarang <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Barang harus dipilih terlebih dahulu (id_barang tidak boleh kosong)",
		})
		return
	}

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

func UpdatePengiriman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method tidak diizinkan",
		})
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ID Pengiriman tidak valid",
		})
		return
	}

	var pengiriman models.Pengiriman
	if err := config.DB.Where("id_pengiriman = ?", id).First(&pengiriman).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Data pengiriman tidak ditemukan",
		})
		return
	}

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