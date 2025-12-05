package services

import (
	"contact-management/src/helpers"
	"contact-management/src/models"
	"contact-management/src/repositories"
	"contact-management/src/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	// "github.com/go-playground/validator/v10"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

var ErrUsernameTaken = errors.New("username sudah digunakan")

var ErrInvalidCredentials = errors.New("username atau password salah")

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (a *AuthService) Register(user *models.User) error {
	validate := helpers.InitValidator()

	err := validate.Struct(user)
	if err != nil {
		formatted := helpers.FormatValidationError(err)
		return helpers.ValidationErrors{Messages: formatted}
	}

	isUser, err := a.userRepo.FindByUsername(user.Username)
	if err != nil && !errors.Is(err, repositories.ErrUserNotFound) {
		return err
	}

	if isUser != nil {
		return ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	err = a.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Login(user *models.User) (string, error) {
	validate := helpers.InitValidator()

	err := validate.Struct(user)
	if err != nil {
		formatted := helpers.FormatValidationError(err)
		return "", helpers.ValidationErrors{Messages: formatted}
	}


	data_user, err := a.userRepo.FindByUsername(user.Username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data_user.Password), []byte(user.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.Username)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) Me(username string) (*models.User, error) {

	user, err := a.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthService) Logout(username, token string) error {
	_, err := a.userRepo.FindByUsername(username)
	if err != nil {
		return err
	}

	err = utils.RevokeToken(token, username)
	if err != nil {
		return err
	}
	return nil
}
