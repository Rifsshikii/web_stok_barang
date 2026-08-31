package main

import (
	"fmt"
	"net/http"
	"web_stok_barang/config"
	"web_stok_barang/controllers"
	"web_stok_barang/middlewares"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	config.ConnectDB()

	http.HandleFunc("/api/register", enableCORS(controllers.Register))
	http.HandleFunc("/api/login", enableCORS(controllers.Login))

	http.HandleFunc("/api/barang", enableCORS(middlewares.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			controllers.GetBarang(w, r)
		} else if r.Method == "POST" {
			controllers.CreateBarang(w, r)
		} else if r.Method == "PUT" {
			controllers.UpdateBarang(w, r)
		} else if r.Method == "DELETE" {
			controllers.DeleteBarang(w, r)
		} else {
			http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		}
	})))

	fmt.Println("Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
