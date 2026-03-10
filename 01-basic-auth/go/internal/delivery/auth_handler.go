package delivery

// HTTP handlers для авторизации (Delivery слой).
//
// Handler'ы — это "входная точка" HTTP:
// - читают и валидируют входные данные (JSON/body)
// - вызывают бизнес-логику (UseCase)
// - формируют HTTP-ответ (status code + JSON)
import (
	"basic-auth/internal/usecase"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

// NewAuthHandler связывает HTTP-слой с UseCase-слоем.
func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Message string      `json:"message"`
	User    interface{} `json:"user,omitempty"`
}

// RegisterHandler — публичный endpoint регистрации.
// При успехе возвращает 201 Created и сообщение.
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authUsecase.Register(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := AuthResponse{
		Message: "Пользователь успешно зарегистрирован",
		User:    user,
	}
	json.NewEncoder(w).Encode(resp)
}

// LoginHandler — простая авторизация по email+password.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := AuthResponse{
		Message: "Пользователь успешно вошёл",
		User:    user,
	}
	json.NewEncoder(w).Encode(resp)
}

// DeleteByCredentialsHandler — удаление пользователя по email+password.
func (h *AuthHandler) DeleteByCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.authUsecase.DeleteByCredentials(req.Email, req.Password); err != nil {
		if err.Error() == "invalid email or password" {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"message": "Пользователь успешно удалён",
	}
	json.NewEncoder(w).Encode(resp)
}
