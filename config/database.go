package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB dipakai untuk controller ber-GORM (seperti pengiriman_controller.go)
var DB *gorm.DB

// DBRAW dipakai untuk controller ber-SQL biasa (seperti auth, barang, barang_keluar)
var DBRAW *sql.DB

func ConnectDB() {
	var err error
	dsn := "root:ikyy1811@tcp(127.0.0.1:3306)/web_stok_barang?charset=utf8mb4&parseTime=True&loc=Local"

	// 1. Inisialisasi GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi GORM:", err)
	}

	// 2. Ekstrak *sql.DB dari GORM untuk SQL Native
	DBRAW, err = DB.DB()
	if err != nil {
		log.Fatal("Gagal mengambil instance sql.DB:", err)
	}

	log.Println("Berhasil terhubung ke database via GORM & SQL Native!")
}
