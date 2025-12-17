package controllers

import (
	"contact-management/src/apps"
	"contact-management/src/helpers"
	"contact-management/src/models"
	"contact-management/src/repositories"
	"contact-management/src/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	apps.LoggingApp().Info("Ini adalah service dari memory", categoryService)
	return &CategoryController{categoryService: categoryService}
}

func (c *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mengambil data kategori", err.Error())
		return
	}
	helpers.SuccessResponse(w, http.StatusOK, "Berhasil mengambil data kategori", categories)
	return
}

func(c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category models.Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		helpers.BadRequestResponse(w, "Gagal memproses input", err)
		return
	}

	err = c.categoryService.CreateCategory(&category)
	if err != nil {
		if validationErr, ok := err.(helpers.ValidationErrors); ok {
			helpers.BadRequestResponse(w, "Gagal membuat kategori", validationErr.Messages)
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal membuat kategori", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusCreated, "Berhasil membuat kategori", category)
	return
}

func (c *CategoryController) GetCategoryByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		helpers.BadRequestResponse(w, "ID harus berupa angka", err)
		return
	}

	category, err := c.categoryService.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorCategoryNotFound) {
			helpers.NotFoundResponse(w, "Kategori tidak ditemukan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mendapatkan data kategori", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil mendapatkan data kategori", category)
	return
}

func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		helpers.BadRequestResponse(w, "ID harus berupa angka", err)
		return
	}

	var category models.Category

	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		helpers.BadRequestResponse(w, "Gagal memproses input", err)
		return
	}

	err = c.categoryService.UpdateCategory(&category, id)
	if err != nil {
		if errors.Is(err, repositories.ErrorCategoryNotFound) {
			helpers.NotFoundResponse(w, "Kategori tidak ditemukan")
			return
		}

		if validationErr, ok := err.(helpers.ValidationErrors); ok {
			helpers.BadRequestResponse(w, "Gagal memperbarui kategori", validationErr.Messages)
			return
		}
		
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal memperbarui kategori", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil memperbarui kategori", category)
	return
}

func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		helpers.BadRequestResponse(w, "ID harus berupa angka", err)
		return
	}

	err = c.categoryService.DeleteCategory(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorCategoryNotFound) {
			helpers.NotFoundResponse(w, "Kategori tidak ditemukan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal menghapus kategori", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil menghapus kategori", nil)
	return
}
