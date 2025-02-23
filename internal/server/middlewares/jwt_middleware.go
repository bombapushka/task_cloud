package middlewares

import (
	"cloud/internal/server/auth"
	"context"
	"net/http"
)

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/login" || r.URL.Path == "/register" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
