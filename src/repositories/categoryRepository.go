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
	UpdateCategory(category *models.Category) error
	DeleteCategory(id int) error
}

func (cr *categoryRepository) CreateCategory(category *models.Category) error {
	result, err := cr.db.Exec("INSERT INTO categories (name) VALUES ($1)", category.Name)
	if err != nil {
		return err
	}

	CategoryID, _ := result.LastInsertId()

	category.CategoryID = int(CategoryID)

	return nil
}

func (cr *categoryRepository) GetAllCategories() ([]*models.Category, error) {
	rows, err := cr.db.Query("SELECT * FROM categories WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		if err := rows.Scan(&category.CategoryID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (cr *categoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	row := cr.db.QueryRow("SELECT * FROM categories WHERE category_id = $1 AND deleted_at IS NULL", id)
	category := &models.Category{}
	if err := row.Scan(&category.CategoryID, &category.Name); err != nil {
		return nil, err
	}
	return category, nil
}

func (cr *categoryRepository) UpdateCategory(category *models.Category) error {
	row, err := cr.db.Exec("UPDATE categories SET name = $2 WHERE category_id = $1 AND deleted_at IS NULL", category.CategoryID, category.Name)
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
	row, err := cr.db.Exec("UPDATE categories SET deleted_at = NOW() WHERE category_id = $1 AND deleted_at IS NULL", id)
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
