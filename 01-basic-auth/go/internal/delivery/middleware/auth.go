package middleware

// Basic Auth middleware (Go).
//
// Что делает:
// - читает заголовок Authorization: Basic <base64(email:password)>
// - вызывает AuthUsecase.Login(email, password)
// - при успехе кладёт user.ID в context, чтобы handler мог взять его без повторной проверки
import (
	"basic-auth/internal/usecase"
	"context"
	"net/http"
)

type contextKey string

// UserIDKey — ключ для хранения user_id в context запроса.
const UserIDKey contextKey = "user_id"

// BasicAuth возвращает middleware-функцию, которую можно "обернуть" вокруг handler'а.
func BasicAuth(authUsecase *usecase.AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// r.BasicAuth() — стандартная функция net/http:
			// парсит заголовок Authorization и возвращает email, password.
			email, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			// Проверяем credentials через бизнес-логику (bcrypt сравнение внутри use case).
			user, err := authUsecase.Login(email, password)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			// Пробрасываем ID пользователя дальше по цепочке обработки.
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
