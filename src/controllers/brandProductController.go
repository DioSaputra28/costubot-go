package controllers

import (
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

type BrandProductController struct {
	brandProductService *services.BrandProductService
}

func NewBrandProductController(brandProductService *services.BrandProductService) *BrandProductController {
	return &BrandProductController{brandProductService: brandProductService}
}

func (bpc *BrandProductController) CreateBrandProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var brandProduct models.BrandProduct

	err := json.NewDecoder(r.Body).Decode(&brandProduct)
	if err != nil {
		helpers.BadRequestResponse(w, "Gagal memproses input", err)
		return
	}

	err = bpc.brandProductService.CreateBrandProduct(&brandProduct)
	if err != nil {
		if validationErr, ok := err.(helpers.ValidationErrors); ok {
			helpers.BadRequestResponse(w, "Gagal membuat brand product", validationErr.Messages)
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal membuat brand product", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusCreated, "Berhasil membuat brand product", brandProduct)
	return
}

func (bpc *BrandProductController) GetAllBrandProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	brandProducts, err := bpc.brandProductService.GetAllBrandProducts()
	if err != nil {
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mendapatkan data brand product", err.Error())
		return
	}
	helpers.SuccessResponse(w, http.StatusOK, "Berhasil mendapatkan data brand product", brandProducts)
	return
}

func (bpc *BrandProductController) GetBrandProductByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		helpers.BadRequestResponse(w, "ID harus berupa angka", err)
		return
	}

	brandProduct, err := bpc.brandProductService.GetBrandProductByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorBrandProductNotFound) {
			helpers.NotFoundResponse(w, "Brand product tidak ditemukan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mendapatkan data brand product", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil mendapatkan data brand product", brandProduct)
	return
}

func (bpc *BrandProductController) UpdateBrandProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		helpers.BadRequestResponse(w, "ID harus berupa angka", err)
		return
	}

	var brandProduct models.BrandProduct

	err = json.NewDecoder(r.Body).Decode(&brandProduct)
	if err != nil {
		helpers.BadRequestResponse(w, "Gagal memproses input", err)
		return
	}

	err = bpc.brandProductService.UpdateBrandProduct(id, &brandProduct)
	if err != nil {
		if errors.Is(err, repositories.ErrorBrandProductNotFound) {
			helpers.NotFoundResponse(w, "Brand product tidak ditemukan")
			return
		}
		if errors.Is(err, repositories.ErrorBrandProductNotFound) {
			helpers.NotFoundResponse(w, "Brand product tidak ditemukan")
			return
		}
		if validationErr, ok := err.(helpers.ValidationErrors); ok {
			helpers.BadRequestResponse(w, "Gagal memperbarui brand product", validationErr.Messages)
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal memperbarui brand product", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil memperbarui brand product", brandProduct)
	return
}

func (bpc *BrandProductController) DeleteBrandProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		helpers.BadRequestResponse(w, "ID harus berupa angka", err)
		return
	}

	err = bpc.brandProductService.DeleteBrandProduct(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorBrandProductNotFound) {
			helpers.NotFoundResponse(w, "Brand product tidak ditemukan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal menghapus brand product", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil menghapus brand product", nil)
	return
}
