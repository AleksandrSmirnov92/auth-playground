package middleware

import (
	"basic-auth/internal/usecase"
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func BasicAuth(authUsecase *usecase.AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			email, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			user, err := authUsecase.Login(email, password)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
