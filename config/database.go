package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	dsn := "root:ikyy1811@tcp(127.0.0.1:3306)/web_stok_barang"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak merespon:", err)
	}
}
