package usecase

// AuthUsecase — бизнес-логика авторизации.
//
// Это слой UseCase: тут нет HTTP и нет деталей хранения данных.
// Он работает через интерфейс domain.UserRepository и решает:
// - можно ли зарегистрировать пользователя (email уникален)
// - как захешировать пароль (bcrypt)
// - как проверить логин/пароль (login)
import (
	"basic-auth/internal/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepository domain.UserRepository
}

// NewAuthUsecase "внедряет" репозиторий в бизнес-логику.
func NewAuthUsecase(userRepository domain.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepository: userRepository}
}

// Register создаёт нового пользователя:
// - проверяет, что email ещё не занят
// - хеширует пароль bcrypt'ом (GenerateFromPassword)
//   (bcrypt берёт «сырой» пароль и возвращает строку-хеш; по хешу нельзя восстановить пароль)
// - сохраняет пользователя в репозиторий
func (u *AuthUsecase) Register(email, password string) (*domain.User, error) {
	existingUser, err := u.userRepository.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	err = u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Login проверяет credentials (email + password).
// Для проверки bcrypt использует CompareHashAndPassword:
//   - на вход: сохранённый хеш (user.Password) и введённый пароль
//   - внутри хешируется введённый пароль тем же алгоритмом и сравнивается с хешем
//   - если хотя бы что-то не совпало — возвращается ошибка.
// Возвращаем одинаковую ошибку для "нет такого email" и "неверный пароль",
// чтобы не помогать атакующему угадывать существующие email (user enumeration).
func (u *AuthUsecase) Login(email, password string) (*domain.User, error) {
	user, err := u.userRepository.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

// GetUserByID возвращает пользователя по ID.
func (u *AuthUsecase) GetUserByID(id string) (*domain.User, error) {
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUserById удаляет пользователя по ID.
// В текущем API вызывается из DeleteByCredentials:
// сначала мы находим пользователя по email+паролю (Login),
// затем удаляем его по ID через этот метод.
func (u *AuthUsecase) DeleteUserById(id string) error {
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteByCredentials удаляет пользователя по email+password.
// Удобно для DELETE /api/v1/auth/delete с телом запроса.
func (u *AuthUsecase) DeleteByCredentials(email, password string) error {
	user, err := u.Login(email, password)
	if err != nil {
		return err
	}
	return u.DeleteUserById(user.ID)
}
