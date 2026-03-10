package delivery

// HTTP handlers для авторизации (Delivery слой).
//
// Handler'ы — это "входная точка" HTTP:
// - читают и валидируют входные данные (JSON/body)
// - вызывают бизнес-логику (UseCase)
// - формируют HTTP-ответ (status code + JSON)
import (
	"basic-auth/internal/delivery/middleware"
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

// RegisterHandler — публичный endpoint регистрации.
// При успехе возвращает 201 Created и пользователя (без password).
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authUsecase.Register(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// MeHandler — защищённый endpoint "кто я".
// user_id берётся из context: туда его положил Basic Auth middleware после успешной проверки.
func (h *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authUsecase.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUserHandler — защищённый endpoint удаления своего аккаунта.
func (h *AuthHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.authUsecase.DeleteUserById(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
