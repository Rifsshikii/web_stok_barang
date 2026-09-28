package models

type Barang struct {
	// Tambahkan gorm:"primaryKey;column:id_barang"
	IDBarang   int     `gorm:"primaryKey;column:id_barang" json:"id_barang"`
	KodeBarang string  `json:"kode_barang"`
	Kategori   string  `json:"kategori"`
	Stok       int     `json:"stok"`
	Harga      float64 `json:"harga"`
	NomorResi  string  `json:"nomor_resi"`
	Status     string  `json:"status"`
}

// Beri tahu GORM nama tabel secara eksplisit
func (Barang) TableName() string {
	return "barang"
}
