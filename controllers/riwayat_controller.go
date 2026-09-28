package controllers

import (
	"encoding/json"
	"net/http"

	"web_stok_barang/config"
)

type RiwayatDTO struct {
	ID         int    `json:"id"`
	Kode       string `json:"kode"`
	KodeBarang string `json:"kode_barang"`
	Tipe       string `json:"tipe"`
	Aktivitas  string `json:"aktivitas"`
	NamaBarang string `json:"nama_barang"`
	Jumlah     int    `json:"jumlah"`
	Tanggal    string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
}

func GetRiwayat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT 
			id, 
			COALESCE(kode, '') AS kode, 
			COALESCE(kode_barang, '') AS kode_barang,
			COALESCE(tipe, '') AS tipe,
			COALESCE(aktivitas, '') AS aktivitas, 
			COALESCE(nama_barang, '') AS nama_barang, 
			COALESCE(jumlah, 0) AS jumlah, 
			COALESCE(DATE_FORMAT(tanggal, '%Y-%m-%d %H:%i:%s'), '') AS tanggal, 
			COALESCE(keterangan, '') AS keterangan 
		FROM riwayat 
		ORDER BY id DESC
	`

	rows, err := config.DBRAW.Query(query)
	if err != nil {
		http.Error(w, "Error DB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	listRiwayat := []RiwayatDTO{}

	for rows.Next() {
		var item RiwayatDTO
		err := rows.Scan(
			&item.ID,
			&item.Kode,
			&item.KodeBarang,
			&item.Tipe,
			&item.Aktivitas,
			&item.NamaBarang,
			&item.Jumlah,
			&item.Tanggal,
			&item.Keterangan,
		)
		if err != nil {
			http.Error(w, "Scan Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		listRiwayat = append(listRiwayat, item)
	}

	json.NewEncoder(w).Encode(listRiwayat)
}
