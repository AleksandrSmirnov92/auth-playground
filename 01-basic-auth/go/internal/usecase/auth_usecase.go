package usecase

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

func NewAuthUsecase(userRepository domain.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepository: userRepository}
}

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

func (u *AuthUsecase) GetUserByID(id string) (*domain.User, error) {
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *AuthUsecase) DeleteUserById(id string) error {
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
