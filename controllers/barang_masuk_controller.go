package controllers

import (
	"encoding/json"
	"net/http"
	"web_stok_barang/config"
)

type BarangMasukInput struct {
	IDBarang   int    `json:"id_barang"`
	Jumlah     int    `json:"jumlah"`
	Tanggal    string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
}

type BarangMasukResponse struct {
	ID         int    `json:"id"`
	IDBarang   int    `json:"id_barang"`
	KodeBarang string `json:"kode_barang"`
	NamaBarang string `json:"nama_barang"`
	Jumlah     int    `json:"jumlah"`
	Tanggal    string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
}

// 1. GET ALL BARANG MASUK
func GetBarangMasuk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// LEFT JOIN + COALESCE untuk cegah scan error pada NULL value
	query := `
		SELECT 
			bm.id, 
			bm.id_barang, 
			COALESCE(b.kode_barang, '-') as kode_barang, 
			COALESCE(b.nama_barang, 'Barang Tidak Ditemukan') as nama_barang, 
			bm.jumlah, 
			DATE_FORMAT(bm.tanggal, '%Y-%m-%d') as tanggal, 
			COALESCE(bm.keterangan, '') as keterangan 
		FROM barang_masuk bm
		LEFT JOIN barang b ON bm.id_barang = b.id_barang
		ORDER BY bm.tanggal DESC, bm.id DESC
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Gagal mengambil data barang masuk: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	listData := make([]BarangMasukResponse, 0)
	for rows.Next() {
		var item BarangMasukResponse
		if err := rows.Scan(
			&item.ID,
			&item.IDBarang,
			&item.KodeBarang,
			&item.NamaBarang,
			&item.Jumlah,
			&item.Tanggal,
			&item.Keterangan,
		); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		listData = append(listData, item)
	}

	json.NewEncoder(w).Encode(listData)
}

// 2. CREATE BARANG MASUK (+ Otomatis Tambah Stok)
func CreateBarangMasuk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input BarangMasukInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	// Validation
	if input.IDBarang == 0 || input.Jumlah <= 0 {
		http.Error(w, "Pilih barang dan masukan jumlah yang valid", http.StatusBadRequest)
		return
	}

	// 1. Simpan ke tabel barang_masuk
	queryInsert := "INSERT INTO barang_masuk (id_barang, jumlah, tanggal, keterangan) VALUES (?, ?, ?, ?)"
	_, err := config.DB.Exec(queryInsert, input.IDBarang, input.Jumlah, input.Tanggal, input.Keterangan)
	if err != nil {
		http.Error(w, "Gagal mencatat barang masuk: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Otomatis Update (tambah) Stok di tabel barang
	queryUpdateStok := "UPDATE barang SET stok = stok + ? WHERE id_barang = ?"
	_, err = config.DB.Exec(queryUpdateStok, input.Jumlah, input.IDBarang)
	if err != nil {
		http.Error(w, "Gagal memperbarui stok barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Barang masuk berhasil dicatat dan stok telah diperbarui!",
	})
}
