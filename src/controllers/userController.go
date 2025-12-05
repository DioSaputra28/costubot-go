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

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userRepo *services.UserService) *UserController {
	return &UserController{
		UserService: userRepo,
	}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query()
	page := 1
	limit := 10
	var err error

	if pageStr := query.Get("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			helpers.BadRequestResponse(w, "Parameter page harus berupa angka positif", err)
			return
		}
	}
	if limitStr := query.Get("per_page"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			helpers.BadRequestResponse(w, "Parameter per_page harus berupa angka positif", err)
			return
		}
	}

	users, err := uc.UserService.GetUsers(page, limit)
	if err != nil {
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mendapatkan data user", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil mendapatkan data user", users)
	return
}
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.BadRequestResponse(w, "Gagal memproses input", err)
		return
	}

	err = uc.UserService.CreateUser(&user)
	if err != nil {
		if validationErr, ok := err.(helpers.ValidationErrors); ok {
			helpers.BadRequestResponse(w, "Gagal membuat user", validationErr.Messages)
			return
		}
		if errors.Is(err, services.ErrUsernameTaken) {
			helpers.BadRequestResponse(w, "Username sudah digunakan", err)
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal membuat user", err.Error())
		return
	}

	result := userResponse{
		UserId: user.UserId,
		Username: user.Username,
		CreatedAt: user.CreatedAt,
	}

	helpers.SuccessResponse(w, http.StatusCreated, "Berhasil membuat user", result)
	return
}


func (uc *UserController) GetUserByUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")

	user, err := uc.UserService.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			helpers.NotFoundResponse(w, "User tidak ditemukan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mendapatkan data user", err.Error())
		return
	}

	result := userResponse{
		UserId: user.UserId,
		Username: user.Username,
		CreatedAt: user.CreatedAt,
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil mendapatkan data user", result)
	return
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.BadRequestResponse(w, "Gagal memproses input", err)
		return
	}

	err = uc.UserService.UpdateUser(username, &user)
	if err != nil {
		if validationErr, ok := err.(helpers.ValidationErrors); ok {
			helpers.BadRequestResponse(w, "Gagal memperbarui user", validationErr.Messages)
			return
		}
		if errors.Is(err, repositories.ErrUserNotFound) {
			helpers.NotFoundResponse(w, "User tidak ditemukan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal memperbarui user", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil memperbarui user", nil)
	return
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")

	err := uc.UserService.DeleteUser(username)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			helpers.NotFoundResponse(w, "User tidak ditemukan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal menghapus user", err.Error())
		return
	}

	helpers.SuccessResponse(w, http.StatusOK, "Berhasil menghapus user", nil)
	return
}