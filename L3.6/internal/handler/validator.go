package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/ginext"
)

// Глобальный экземпляр валидатора
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError представляет ошибку валидации одного поля
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors форматирует ошибки валидации в читаемый формат
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			var message string

			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("Поле '%s' обязательно для заполнения", e.Field())
			case "oneof":
				message = fmt.Sprintf("Поле '%s' должно быть одним из: %s", e.Field(), e.Param())
			case "gt":
				message = fmt.Sprintf("Поле '%s' должно быть больше %s", e.Field(), e.Param())
			case "len":
				message = fmt.Sprintf("Поле '%s' должно содержать %s символов", e.Field(), e.Param())
			case "max":
				message = fmt.Sprintf("Поле '%s' не должно превышать %s символов", e.Field(), e.Param())
			case "uuid":
				message = fmt.Sprintf("Поле '%s' должно быть валидным UUID", e.Field())
			default:
				message = fmt.Sprintf("Поле '%s' не прошло валидацию '%s'", e.Field(), e.Tag())
			}

			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: message,
			})
		}
	}

	return errors
}

// ValidateStruct валидирует структуру и возвращает отформатированные ошибки
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// RespondWithValidationError отправляет ответ с ошибками валидации
func RespondWithValidationError(c *ginext.Context, err error) {
	errors := FormatValidationErrors(err)
	c.JSON(http.StatusNotFound, ginext.H{
		"error":  "Ошибка валидации",
		"fields": errors,
	})
}
