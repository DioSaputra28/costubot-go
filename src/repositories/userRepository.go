package repositories

import (
	"contact-management/src/models"
	"database/sql"
)

type UserRepository interface {
	FindByUsername(username string) (any, error)
	CreateUser(user *models.User) error
	GetUser(any) (any, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(user *models.User) error {
	result, err := u.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)	
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.UserId = int(id)

	return nil
}

func (u *userRepository) FindByUsername(username string) (any, error) {
	var user models.User
	err := u.db.QueryRow("SELECT user_id, username, password FROM users WHERE username = ?", username).Scan(&user.UserId, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}