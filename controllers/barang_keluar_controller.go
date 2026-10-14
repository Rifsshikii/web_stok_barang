package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"web_stok_barang/config"
)

type BarangKeluarInput struct {
	IDBarang     int    `json:"id_barang"`
	Jumlah       int    `json:"jumlah"`
	Tanggal      string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
}

type BarangKeluarResponse struct {
	ID           int    `json:"id"`
	IDBarang     int    `json:"id_barang"`
	KodeBarang string `json:"kode_barang"`
	NamaBarang string `json:"nama_barang"`
	Kategori     string `json:"kategori"`
	Jumlah       int    `json:"jumlah"`
	Tanggal      string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
}

func GetBarangKeluar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT 
			bk.id, 
			bk.id_barang, 
			COALESCE(b.kode_barang, '-') as kode_barang, 
			COALESCE(b.nama_barang, 'Barang Tidak Ditemukan') as nama_barang, 
			COALESCE(b.kategori, '-') as kategori, 
			bk.jumlah, 
			DATE_FORMAT(bk.tanggal, '%Y-%m-%d') as tanggal, 
			COALESCE(bk.keterangan, '') as keterangan 
		FROM barang_keluar bk
		LEFT JOIN barang b ON bk.id_barang = b.id_barang
		ORDER BY bk.tanggal DESC, bk.id DESC
	`
	rows, err := config.DBRAW.Query(query)
	if err != nil {
		http.Error(w, "Gagal mengambil data barang keluar: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	listData := make([]BarangKeluarResponse, 0)
	for rows.Next() {
		var item BarangKeluarResponse
		if err := rows.Scan(
			&item.ID,
			&item.IDBarang,
			&item.KodeBarang,
			&item.NamaBarang,
			&item.Kategori,
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

func CreateBarangKeluar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input BarangKeluarInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Format data tidak valid", http.StatusBadRequest)
		return
	}

	if input.IDBarang == 0 || input.Jumlah <= 0 {
		http.Error(w, "Pilih barang dan masukkan jumlah yang valid", http.StatusBadRequest)
		return
	}

	var stokSaatIni int
	err := config.DBRAW.QueryRow("SELECT stok FROM barang WHERE id_barang = ?", input.IDBarang).Scan(&stokSaatIni)
	if err == sql.ErrNoRows {
		http.Error(w, "Barang tidak ditemukan di database", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Gagal memeriksa stok barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if stokSaatIni < input.Jumlah {
		http.Error(w, "Stok barang tidak mencukupi!", http.StatusBadRequest)
		return
	}

	queryInsert := "INSERT INTO barang_keluar (id_barang, jumlah, tanggal, keterangan) VALUES (?, ?, ?, ?)"
	_, err = config.DBRAW.Exec(queryInsert, input.IDBarang, input.Jumlah, input.Tanggal, input.Keterangan)
	if err != nil {
		http.Error(w, "Gagal mencatat barang keluar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	queryUpdateStok := "UPDATE barang SET stok = stok - ? WHERE id_barang = ?"
	_, err = config.DBRAW.Exec(queryUpdateStok, input.Jumlah, input.IDBarang)
	if err != nil {
		http.Error(w, "Gagal memperbarui stok barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Barang keluar berhasil dicatat dan stok telah dikurangi!",
	})
}