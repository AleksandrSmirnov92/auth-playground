package domain

import "time"

// User — доменная сущность пользователя.
//
// Зачем она нужна:
// - Repository хранит/ищет пользователей (Create/GetByID/GetByEmail/Delete)
// - UseCase создаёт пользователя при регистрации и возвращает его в Delivery
// - Delivery сериализует пользователя в JSON-ответ (кроме пароля)
//
// Про теги `json:"..."`:
// - Это struct tags: подсказка пакету encoding/json, как назвать поле в JSON.
// - `json:"id"` означает: в JSON поле будет называться "id" (а не "ID").
// - `json:"created_at"` — snake_case, чтобы API было удобнее и единообразнее.
// - `json:"-"` означает: поле НЕ включать в JSON (пароль никогда не отдаём наружу).
type User struct {
	// ID — уникальный идентификатор пользователя (UUID в текущей реализации).
	ID string `json:"id"`

	// Name — имя пользователя. В этом проекте сейчас не заполняется, оставлено для примера.
	Name string `json:"name"`

	// Email — логин пользователя. Используется для регистрации и для Basic Auth (email:password).
	Email string `json:"email"`

	// Password — bcrypt-хеш пароля. Тег `json:"-"` гарантирует, что поле не уйдёт в JSON-ответ.
	Password string `json:"-"`

	// CreatedAt — время регистрации.
	CreatedAt time.Time `json:"created_at"`
}
