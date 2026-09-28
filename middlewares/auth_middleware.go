package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Gunakan Huruf Kapital (JWTSecret) agar bisa diimpor dan dipakai oleh auth_controller.go
var JWTSecret = []byte("rifsshikii1811")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token tidak ditemukan", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
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
