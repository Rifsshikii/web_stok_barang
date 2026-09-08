package main

import (
	"fmt"
	"net/http"
	"strings"

	"web_stok_barang/config"
	"web_stok_barang/controllers"
	"web_stok_barang/middlewares"
	"web_stok_barang/policy"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Jika HTTP Request method OPTIONS, langsung return OK tanpa lewat auth middleware
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper untuk melewatkan OPTIONS sebelum mengeksekusi middleware
func handleProtected(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func main() {
	config.ConnectDB()

	mux := http.NewServeMux()

	// 1. Route Public
	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)

	// 2. Handler Pengaturan
	mux.HandleFunc("/api/pengaturan", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.GetPengaturan(w, r)
				case http.MethodPut, http.MethodPost:
					middlewares.PolicyMiddleware(policy.RoleAdmin)(controllers.UpdatePengaturan)(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	// 3. Handler Barang Masuk
	mux.HandleFunc("/api/barang-masuk", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.GetBarangMasuk(w, r)
				case http.MethodPost:
					middlewares.PolicyMiddleware(policy.RoleAdmin)(controllers.CreateBarangMasuk)(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	// 4. Handler Barang Keluar
	mux.HandleFunc("/api/barang-keluar", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.GetBarangKeluar(w, r)
				case http.MethodPost:
					middlewares.PolicyMiddleware(policy.RoleAdmin)(controllers.CreateBarangKeluar)(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	// 5. Handler Riwayat Transaksi
	mux.HandleFunc("/api/riwayat", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					controllers.GetRiwayat(w, r)
				} else {
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	// 6. Handler Data User (Khusus Admin)
	mux.HandleFunc("/api/users", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin)(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					controllers.GetUsers(w, r)
				} else {
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	// 7. Handler /api/barang dan /api/barang/{id}
	mux.HandleFunc("/api/barang", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.GetBarang(w, r)
				case http.MethodPost:
					middlewares.PolicyMiddleware(policy.RoleAdmin)(controllers.CreateBarang)(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	mux.HandleFunc("/api/barang/", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin)(
			func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/api/barang-masuk") ||
					strings.HasPrefix(r.URL.Path, "/api/barang-keluar") ||
					strings.HasPrefix(r.URL.Path, "/api/riwayat") ||
					strings.HasPrefix(r.URL.Path, "/api/users") ||
					strings.HasPrefix(r.URL.Path, "/api/pengaturan") {
					http.NotFound(w, r)
					return
				}

				switch r.Method {
				case http.MethodPut:
					controllers.UpdateBarang(w, r)
				case http.MethodDelete:
					controllers.DeleteBarang(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	fmt.Println("Server berjalan di http://localhost:8080")

	err := http.ListenAndServe(":8080", enableCORS(mux))
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
