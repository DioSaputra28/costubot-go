package controllers

import (
	"contact-management/src/helpers"
	"contact-management/src/models"
	"contact-management/src/services"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)


type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (a *AuthController) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		helpers.ErrorResponse(w, http.StatusBadRequest, "Input tidak boleh kosong", "request body is nil")
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.ErrorResponse(w, http.StatusBadRequest, "Gagal memproses input", err.Error())
		return
	}

	err = a.AuthService.Register(&user)
	if err != nil {
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mendaftar user", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusCreated, "Berhasil mendaftar user", nil)
}