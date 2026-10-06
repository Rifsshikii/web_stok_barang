package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// GetJWTSecret mengambil kunci rahasia dari variabel lingkungan (.env)
// Bisa dipanggil oleh auth_controller.go saat membuat (sign) token
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "20091118" // Fallback jika .env belum sempat dibaca
	}
	return []byte(secret)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token tidak ditemukan", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token tidak valid atau expired", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Gagal memproses token claims", http.StatusUnauthorized)
			return
		}

		// Simpan role & id ke context agar bisa dibaca PolicyMiddleware atau Controller
		ctx := context.WithValue(r.Context(), "role", claims["role"])
		ctx = context.WithValue(ctx, "id", claims["id"])

		next(w, r.WithContext(ctx))
	}
}