package repositories

import (
	"contact-management/src/models"
	"database/sql"
	"errors"
)

type brandProductRepository struct {
	db *sql.DB
}

func NewBrandProductRepository(db *sql.DB) *brandProductRepository {
	return &brandProductRepository{db: db}
}

var ErrorBrandProductNotFound = errors.New("brand product not found")

type BrandProductRepository interface {
	CreateBrandProduct(brandProduct *models.BrandProduct) error
	GetAllBrandProducts() ([]models.BrandProduct, error)
	GetBrandProductByID(id int) (*models.BrandProduct, error)
	UpdateBrandProduct(id int, brandProduct *models.BrandProduct) error
	DeleteBrandProduct(id int) error
	GetCategoryByID(id int) (*models.Category, error)
}

func (bpr *brandProductRepository) CreateBrandProduct(brandProduct *models.BrandProduct) error {
	result, err := bpr.db.Exec("INSERT INTO brand_products (name, category_id) VALUES (?, ?)", brandProduct.Name, brandProduct.CategoryID)
	if err != nil {
		return err
	}

	brandProductID, _ := result.LastInsertId()
	brandProduct.BrandProductID = int(brandProductID)
	return nil
}

func (bpr *brandProductRepository) GetAllBrandProducts() ([]models.BrandProduct, error) {
	rows, err := bpr.db.Query("SELECT brand_product_id, name, category_id, created_at, updated_at, deleted_at FROM brand_products WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brandProducts []models.BrandProduct
	for rows.Next() {
		brandProduct := models.BrandProduct{}
		var deletedAt sql.NullTime
		if err := rows.Scan(&brandProduct.BrandProductID, &brandProduct.Name, &brandProduct.CategoryID, &brandProduct.CreatedAt, &brandProduct.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			brandProduct.DeletedAt = &deletedAt.Time
		}
		brandProducts = append(brandProducts, brandProduct)
	}

	return brandProducts, nil
}

func (bpr *brandProductRepository) GetBrandProductByID(id int) (*models.BrandProduct, error) {
	row := bpr.db.QueryRow("SELECT brand_product_id, name, category_id, created_at, updated_at, deleted_at FROM brand_products WHERE brand_product_id = ? AND deleted_at IS NULL", id)
	brandProduct := models.BrandProduct{}
	var deletedAt sql.NullTime
	if err := row.Scan(&brandProduct.BrandProductID, &brandProduct.Name, &brandProduct.CategoryID, &brandProduct.CreatedAt, &brandProduct.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorBrandProductNotFound
		}
		return nil, err
	}
	if deletedAt.Valid {
		brandProduct.DeletedAt = &deletedAt.Time
	}
	return &brandProduct, nil
}

func (bpr *brandProductRepository) UpdateBrandProduct(id int, brandProduct *models.BrandProduct) error {
	row, err := bpr.db.Exec("UPDATE brand_products SET name = ?, category_id = ? WHERE brand_product_id = ? AND deleted_at IS NULL", brandProduct.Name, brandProduct.CategoryID, id)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrorBrandProductNotFound
	}
	return nil
}

func (bpr *brandProductRepository) DeleteBrandProduct(id int) error {
	row, err := bpr.db.Exec("UPDATE brand_products SET deleted_at = NOW() WHERE brand_product_id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrorBrandProductNotFound
	}
	return nil
}

func (bpr *brandProductRepository) GetCategoryByID(id int) (*models.Category, error) {
	row := bpr.db.QueryRow("SELECT category_id, name, created_at, updated_at, deleted_at FROM category WHERE category_id = ? AND deleted_at IS NULL", id)
	category := models.Category{}
	var deletedAt sql.NullTime
	if err := row.Scan(&category.CategoryID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &deletedAt); err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		category.DeletedAt = &deletedAt.Time
	}
	return &category, nil
}
