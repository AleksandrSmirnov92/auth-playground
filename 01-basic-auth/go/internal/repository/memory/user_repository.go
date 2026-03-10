package memory

// In-memory реализация репозитория пользователей.
//
// Важно:
// - данные хранятся в map в памяти процесса (после перезапуска сервера всё исчезнет)
// - RWMutex защищает map при одновременных запросах (конкурентный доступ)
import (
	"basic-auth/internal/domain"
	"errors"
	"sync"
)

type UserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

// NewUserRepository создаёт пустое in-memory хранилище.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

// Create сохраняет нового пользователя (по ключу user.ID).
func (r *UserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

// GetByID ищет пользователя по ID (O(1) по map).
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetByEmail ищет пользователя по email.
// В in-memory варианте это перебор всех пользователей (O(n)).
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Delete удаляет пользователя по ID.
func (r *UserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, id)
	return nil
}
