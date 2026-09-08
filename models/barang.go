package models

type Barang struct {
	IDBarang   int     `json:"id_barang"`
	KodeBarang string  `json:"kode_barang"`
	Kategori   string  `json:"kategori"`
	Stok       int     `json:"stok"`
	Harga      float64 `json:"harga"`
	NomorResi  string  `json:"nomor_resi"`
	Status     string  `json:"status"`
}
