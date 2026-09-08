package middlewares

import (
	"net/http"
)

func PolicyMiddleware(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("role").(string)
			if !ok || userRole == "" {
				http.Error(w, "Forbidden: Role tidak ditemukan", http.StatusForbidden)
				return
			}

			isAllowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Forbidden: Anda tidak memiliki akses ke endpoint ini", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}
