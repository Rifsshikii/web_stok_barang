package models

type Barang struct {
	// Tambahkan gorm:"primaryKey;column:id_barang"
	IDBarang   int     `gorm:"primaryKey;column:id_barang" json:"id_barang"`
	KodeBarang string  `gorm:"column:kode_barang" json:"kode_barang"`
	NamaBarang string  `gorm:"column:nama_barang" json:"nama_barang"` // <--- TAMBAHKAN BARIS INI
	Kategori   string  `gorm:"column:kategori" json:"kategori"`
	Stok       int     `gorm:"column:stok" json:"stok"`
	Harga      float64 `gorm:"column:harga" json:"harga"`
	NomorResi  string  `gorm:"column:nomor_resi" json:"nomor_resi"`
	Status     string  `gorm:"column:status" json:"status"`
}

// Beri tahu GORM nama tabel secara eksplisit
func (Barang) TableName() string {
	return "barang"
}