package repositories

import (
	"contact-management/src/models"
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("User tidak ditemukan")

type UserRepository interface {
	GetUsers(page, limit int) ([]models.User, int, error)
	FindByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(username string, user *models.User) (int64, error)
	DeleteUser(username string) (int64, error)
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

func (u *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow("SELECT user_id, username, password FROM users WHERE username = ?", username).Scan(&user.UserId, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetUsers(page, limit int) ([]models.User, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	rows, err := u.db.Query("SELECT user_id, username, created_at, updated_at FROM users LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]models.User, 0, limit)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserId, &user.Username, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	var total int
	err = u.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *userRepository) UpdateUser(username string, user *models.User) (int64, error) {
	result, err := u.db.Exec("UPDATE users SET username = ?, password = ? WHERE username = ?", user.Username, user.Password, username)
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	user.UserId = int(id)

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return RowsAffected, err
}

func (u *userRepository) DeleteUser(username string) (int64, error)  {
	row, err := u.db.Exec("DELETE FROM users WHERE username = ?", username)
	RowsAffected, err := row.RowsAffected()
	if err != nil {
		return 0, err
	}
	return RowsAffected, err
}