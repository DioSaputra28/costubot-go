package services

import (
	"contact-management/src/helpers"
	"contact-management/src/models"
	"contact-management/src/repositories"
)

type BrandProductService struct {
	brandProductRepository repositories.BrandProductRepository
}

func NewBrandProductService(brandProductRepository repositories.BrandProductRepository) *BrandProductService {
	return &BrandProductService{brandProductRepository: brandProductRepository}
}

func (bps *BrandProductService) CreateBrandProduct(brandProduct *models.BrandProduct) error {
	validate := helpers.InitValidator()

	err := validate.Struct(brandProduct)
	if err != nil {
		formatted := helpers.FormatValidationError(err)
		return helpers.ValidationErrors{Messages: formatted}
	}

	category, err := bps.brandProductRepository.GetCategoryByID(brandProduct.CategoryID)
	if err != nil {
		return err
	}

	if category == nil {
		return repositories.ErrorCategoryNotFound
	}

	return bps.brandProductRepository.CreateBrandProduct(brandProduct)
}

func (bps *BrandProductService) GetAllBrandProducts() ([]models.BrandProduct, error) {
	return bps.brandProductRepository.GetAllBrandProducts()
}

func (bps *BrandProductService) GetBrandProductByID(id int) (*models.BrandProduct, error) {
	brandProduct, err := bps.brandProductRepository.GetBrandProductByID(id)
	if err != nil {
		return nil, err
	}

	if brandProduct == nil {
		return nil, repositories.ErrorBrandProductNotFound
	}

	return brandProduct, nil
}

func (bps *BrandProductService) UpdateBrandProduct(id int, brandProduct *models.BrandProduct) error {
	validate := helpers.InitValidator()

	err := validate.Struct(brandProduct)
	if err != nil {
		formatted := helpers.FormatValidationError(err)
		return helpers.ValidationErrors{Messages: formatted}
	}

	category, err := bps.brandProductRepository.GetCategoryByID(brandProduct.CategoryID)
	if err != nil {
		return err
	}

	if category == nil {
		return repositories.ErrorCategoryNotFound
	}

	return bps.brandProductRepository.UpdateBrandProduct(id, brandProduct)
}

func (bps *BrandProductService) DeleteBrandProduct(id int) error {
	brandProduct, err := bps.brandProductRepository.GetBrandProductByID(id)
	if err != nil {
		return err
	}

	if brandProduct == nil {
		return repositories.ErrorBrandProductNotFound
	}

	return bps.brandProductRepository.DeleteBrandProduct(id)
}
