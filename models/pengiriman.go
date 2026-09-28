package models

import "time"

type Pengiriman struct {
	IDPengiriman      int       `gorm:"primaryKey;column:id_pengiriman" json:"id_pengiriman"`
	IDBarang          int       `gorm:"column:id_barang" json:"id_barang"`
	Barang            Barang    `gorm:"foreignKey:IDBarang;references:IDBarang" json:"barang"` // Relasi ke struct Barang
	Jumlah            int       `gorm:"column:jumlah" json:"jumlah"`
	Tujuan            string    `gorm:"column:tujuan" json:"tujuan"`
	TanggalPengiriman string    `gorm:"column:tanggal_pengiriman" json:"tanggal_pengiriman"`
	Keterangan        string    `gorm:"column:keterangan" json:"keterangan"`
	Status            string    `gorm:"column:status" json:"status"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Pengiriman) TableName() string {
	return "pengiriman"
}
