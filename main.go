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
	config.ConnectDB()

	mux := http.NewServeMux()

	// 1. Route Public
	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)

	// 2. Profile / Check Session (Dapat diakses Admin, Petugas, dan User)
	mux.HandleFunc("/api/me", handleProtected(middlewares.AuthMiddleware(controllers.GetProfile)))

	// 3. Handler Pengaturan (GET: Admin & Petugas, POST/PUT: Admin)
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

	// 4. Handler Barang Masuk (POST sekarang bisa diakses Petugas juga)
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

	// 5. Handler Barang Keluar (POST sekarang bisa diakses Petugas juga)
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

	// 6. Handler Riwayat Transaksi (Diakses Admin, Petugas, dan User/Pengguna)
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

	// 7. Handler Data User (Khusus Admin)
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

	// 7b. Handler Hapus User berdasarkan ID (/api/users/id) (Khusus Admin)
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

	// 8. Handler Pengiriman Barang (GET & POST) (Petugas & Admin)
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

	// 8b. TAMBAHAN BARU: Handler Update Status Pengiriman berdasarkan ID (/api/pengiriman/id) (Petugas & Admin)
	mux.HandleFunc("/api/pengiriman/", handleProtected(middlewares.AuthMiddleware(
		middlewares.PolicyMiddleware(policy.RoleAdmin, policy.RolePetugas)(
			func(w http.ResponseWriter, r *http.Request) {
				// Abaikan jika rute ini sebenarnya memanggil endpoint lain yang kebetulan berawalan sama
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

	// 9. Handler /api/barang (GET dapat diakses oleh Admin, Petugas, dan User)
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

	fmt.Println("Server berjalan di http://localhost:8080")

	err := http.ListenAndServe(":8080", enableCORS(mux))
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
