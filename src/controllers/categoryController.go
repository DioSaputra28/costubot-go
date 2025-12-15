package controllers

import (
	"contact-management/src/apps"
	"contact-management/src/services"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	apps.LoggingApp().Info("Ini adalah service dari memory", categoryService)
	return &CategoryController{categoryService: categoryService}
}

func(c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
}