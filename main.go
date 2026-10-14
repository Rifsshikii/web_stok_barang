package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variabel lingkungan bawaan")
	}

	config.ConnectDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)

	mux.HandleFunc("/api/me", handleProtected(middlewares.AuthMiddleware(controllers.GetProfile)))

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

	mux.HandleFunc("/api/barang-masuk", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.GetBarangMasuk(w, r)
				case http.MethodPost:
					middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(controllers.CreateBarangMasuk)(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	mux.HandleFunc("/api/barang-keluar", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.GetBarangKeluar(w, r)
				case http.MethodPost:
					middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(controllers.CreateBarangKeluar)(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	mux.HandleFunc("/api/riwayat", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas, policy.RoleUser)(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					controllers.GetRiwayat(w, r)
				} else {
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

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

	mux.HandleFunc("/api/users/", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodDelete:
					controllers.DeleteUser(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	mux.HandleFunc("/api/pengiriman", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.GetPengiriman(w, r)
				case http.MethodPost:
					controllers.CreatePengiriman(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	mux.HandleFunc("/api/pengiriman/", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/api/pengaturan") {
					http.NotFound(w, r)
					return
				}

				switch r.Method {
				case http.MethodPut:
					controllers.UpdatePengiriman(w, r)
				default:
					http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
				}
			},
		),
	)))

	mux.HandleFunc("/api/barang", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas, policy.RoleUser)(
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
					strings.HasPrefix(r.URL.Path, "/api/pengiriman") ||
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server berjalan di http://localhost:%s\n", port)

	err := http.ListenAndServe(":"+port, enableCORS(mux))
	if err != nil {
		fmt.Println("Server error:", err)
	}
}