package services

import (
	"contact-management/src/helpers"
	"contact-management/src/models"
	"contact-management/src/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (us *UserService) GetUsers(page, limit int) (*models.UserResponsePagination, error) {
	users, total, err := us.userRepo.GetUsers(page, limit)
	if err != nil {
		return nil, err
	}

	responseUsers := make([]models.UserResponse, len(users))
	for i, user := range users {
		responseUsers[i] = models.UserResponse{
			UserId:    user.UserId,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		}
	}

	return &models.UserResponsePagination{
		Users: responseUsers,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (uc *UserService) CreateUser(user *models.User) error {

	validate := helpers.InitValidator()

	err := validate.Struct(user)
	if err != nil {
		formatted := helpers.FormatValidationError(err)
		return helpers.ValidationErrors{Messages: formatted}
	}
	
	dataUser, _ := uc.userRepo.FindByUsername(user.Username)
	if dataUser != nil {
		return ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = uc.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserService) UpdateUser(username string, user *models.User) error {
	validate := helpers.InitValidator()

	err := validate.Struct(user)
	if err != nil {
		formatted := helpers.FormatValidationError(err)
		return helpers.ValidationErrors{Messages: formatted}
	}

	dataUser, _ := uc.userRepo.FindByUsername(user.Username)
	if dataUser != nil {
		return ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	row, err := uc.userRepo.UpdateUser(username ,user)
	if err != nil {
		return err
	}
	
	if row == 0 {
		return repositories.ErrUserNotFound
	}

	return nil
}

func (uc *UserService) DeleteUser(username string) error {
	row, err := uc.userRepo.DeleteUser(username)
	if err != nil {
		return err
	}

	if row == 0 {
		return repositories.ErrUserNotFound
	}

	return nil
}