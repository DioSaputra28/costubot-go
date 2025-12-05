package helpers

import (
	"strings"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

var (
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate *validator.Validate
)

// InitValidator inisialisasi validator dengan bahasa Indonesia
func InitValidator() *validator.Validate {
	// Setup translator bahasa Indonesia
	indonesian := id.New()
	uni = ut.New(indonesian, indonesian)
	trans, _ = uni.GetTranslator("id")

	// Buat validator instance
	validate = validator.New()

	// Register default translation bahasa Indonesia
	id_translations.RegisterDefaultTranslations(validate, trans)

	// Custom translation untuk field tertentu
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} wajib diisi", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} minimal {1} karakter", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	validate.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0} maksimal {1} karakter", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field(), fe.Param())
		return t
	})

	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} harus format email yang valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	return validate
}

func FormatValidationError(err error) map[string]string {
	errorMessages := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		translatedErrors := validationErrors.Translate(trans)

		for field, message := range translatedErrors {
			fieldName := strings.ToLower(field)
			errorMessages[fieldName] = message
		}
	}

	return errorMessages
}

type ValidationErrors struct {
	Messages map[string]string
}

func (v ValidationErrors) Error() string {
	return "validation error"
}
