package services

import (
	"contact-management/src/helpers"
	"contact-management/src/models"
	"contact-management/src/repositories"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (cs *CategoryService) GetAllCategories() ([]*models.Category, error) {
	categories, err := cs.categoryRepo.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (cs *CategoryService) GetCategoryByID(id int) (*models.Category, error) {

	return cs.categoryRepo.GetCategoryByID(id)
}

func (cs *CategoryService) CreateCategory(category *models.Category) error {
	validate := helpers.InitValidator()

	err := validate.Struct(category)
	if err != nil {
		return err
	}
	return cs.categoryRepo.CreateCategory(category)
}

func (cs *CategoryService) UpdateCategory(category *models.Category, id int) error {
	validate := helpers.InitValidator()

	err := validate.Struct(category)
	if err != nil {
		return err
	}
	return cs.categoryRepo.UpdateCategory(category, id)
}

func (cs *CategoryService) DeleteCategory(id int) error {
	return cs.categoryRepo.DeleteCategory(id)
}
