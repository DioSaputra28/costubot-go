package services

import (
	"contact-management/src/models"
	"contact-management/src/repositories"

	"github.com/go-playground/validator/v10"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (a *AuthService) Register(user *models.User) error {
	validate := validator.New()

	err := validate.Struct(user)
	if err != nil {
		return err
	}

	_, err = a.userRepo.FindByUsername(user.Username)
	if err != nil {
		return err
	}

	err = a.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}