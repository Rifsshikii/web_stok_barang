package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	// 1. Ambil kredensial database dari file .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Nilai default jika .env belum sempat terbaca
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "root"
	}
	if dbName == "" {
		dbName = "web_stok_barang"
	}

	// 2. Susun DSN secara dinamis
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// 3. Inisialisasi GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi GORM:", err)
	}

	// 4. Ekstrak *sql.DB dari GORM untuk SQL Native
	DBRAW, err = DB.DB()
	if err != nil {
		log.Fatal("Gagal mengambil instance sql.DB:", err)
	}

	log.Println("Berhasil terhubung ke database via GORM & SQL Native!")
}