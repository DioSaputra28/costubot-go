package repositories

import (
	"contact-management/src/models"
	"database/sql"
	"errors"
)

var ErrorCategoryNotFound = errors.New("category not found")

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *categoryRepository {
	return &categoryRepository{db}
}

type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetAllCategories() ([]*models.Category, error)
	GetCategoryByID(id int) (*models.Category, error)
	UpdateCategory(category *models.Category, id int) error
	DeleteCategory(id int) error
}

func (cr *categoryRepository) CreateCategory(category *models.Category) error {
	result, err := cr.db.Exec("INSERT INTO category (name) VALUES (?)", category.Name)
	if err != nil {
		return err
	}

	CategoryID, _ := result.LastInsertId()

	category.CategoryID = int(CategoryID)

	return nil
}

func (cr *categoryRepository) GetAllCategories() ([]*models.Category, error) {
	rows, err := cr.db.Query("SELECT category_id, name, created_at, updated_at, deleted_at FROM category WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		var deletedAt sql.NullTime
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			category.DeletedAt = &deletedAt.Time
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (cr *categoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	row := cr.db.QueryRow("SELECT category_id, name, created_at, updated_at, deleted_at FROM category WHERE category_id = ? AND deleted_at IS NULL", id)
	category := &models.Category{}
	var deletedAt sql.NullTime
	if err := row.Scan(&category.CategoryID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &deletedAt); err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		category.DeletedAt = &deletedAt.Time
	}
	return category, nil
}

func (cr *categoryRepository) UpdateCategory(category *models.Category, id int) error {
	row, err := cr.db.Exec("UPDATE category SET name = ? WHERE category_id = ? AND deleted_at IS NULL", category.Name, id)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrorCategoryNotFound
	}
	return nil
}

func (cr *categoryRepository) DeleteCategory(id int) error {
	row, err := cr.db.Exec("UPDATE category SET deleted_at = NOW() WHERE category_id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrorCategoryNotFound
	}
	return nil
}
