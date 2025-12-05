package controllers

import (
	"contact-management/src/helpers"
	"contact-management/src/models"
	"contact-management/src/services"
	"time"

	// "contact-management/src/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AuthController struct {
	AuthService *services.AuthService
}

type userResponse struct {
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (a *AuthController) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.ErrorResponse(w, http.StatusBadRequest, "Gagal memproses input", err.Error())
		return
	}

	err = a.AuthService.Register(&user)
	if err != nil {
		if validationErrs, ok := err.(helpers.ValidationErrors); ok {
			helpers.ValidationErrorResponse(w, "Validasi gagal", validationErrs.Messages)
			return
		}
		if errors.Is(err, services.ErrUsernameTaken) {
			helpers.ConflictResponse(w, "Username sudah digunakan")
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal mendaftar user", err.Error())
		return
	}

	result := userResponse{
		UserId:    user.UserId,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	helpers.SuccessResponse(w, http.StatusCreated, "Berhasil mendaftar user", result)
	return
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user *models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.ErrorResponse(w, http.StatusBadRequest, "Gagal memproses input", err.Error())
		return
	}

	token, err := a.AuthService.Login(user)
	if err != nil {
		if validationErrs, ok := err.(helpers.ValidationErrors); ok {
			helpers.ValidationErrorResponse(w, "Validasi gagal", validationErrs.Messages)
			return
		}
		if errors.Is(err, services.ErrInvalidCredentials) {
			helpers.ErrorResponse(w, http.StatusUnauthorized, "Username atau password salah", err.Error())
			return
		}
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal login", err.Error())
		return
	}
	helpers.SuccessResponse(w, http.StatusOK, "Berhasil login", map[string]string{"token": token})
	return
}

func (a *AuthController) Me(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {

	username := r.Context().Value("username").(string)
	user, err := a.AuthService.Me(username)
	if err != nil {
		helpers.ErrorResponse(w, http.StatusUnauthorized, "Gagal mengambil informasi user", err.Error())
		return
	}
	helpers.SuccessResponse(w, http.StatusOK, "Berhasil mengambil informasi user", userResponse{
		UserId:    user.UserId,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Context().Value("username").(string)
	err := a.AuthService.Logout(username, r.Header.Get("Authorization"))
	if err != nil {
		helpers.ErrorResponse(w, http.StatusInternalServerError, "Gagal logout", err.Error())
		return
	}
	helpers.SuccessResponse(w, http.StatusOK, "Berhasil logout", nil)
}
