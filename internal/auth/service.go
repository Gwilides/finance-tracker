package auth

import (
	"errors"

	"github.com/Gwilides/finance-tracker/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type UserProvider interface {
	Create(user *user.User) error
	FindByEmail(email string) (*user.User, error)
}

type AuthService struct {
	repository UserProvider
}

func NewAuthService(repository UserProvider) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (service *AuthService) Login(req *LoginRequest) (string, error) {
	existedUser, err := service.repository.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return existedUser.Email, nil
}

func (service *AuthService) Register(req *RegisterRequest) (string, error) {
	_, err := service.repository.FindByEmail(req.Email)
	if err == nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	err = service.repository.Create(&user.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	})
	if err != nil {
		return "", err
	}
	return req.Email, nil
}
