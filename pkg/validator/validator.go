package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Здесь регистрация кастомных валидаций
	// Например: v.RegisterValidation("strong_password", validateStrongPassword)

	return &Validator{validate: v}
}

func (v *Validator) ValidateStruct(s any) error {
	return v.validate.Struct(s)
}

func FormatErrors(err error) []ValidationErrorResponse {
	var errors []ValidationErrorResponse

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			var msg string

			switch err.Tag() {
			// --- Базовые проверки на наличие ---
			case "required":
				msg = fmt.Sprintf("%s is a required field", err.Field())
			case "required_if", "required_unless", "required_with", "required_without":
				// err.Param() здесь будет содержать имя другого поля, от которого зависит наличие текущего
				msg = fmt.Sprintf("%s is required based on the condition of %s", err.Field(), err.Param())

			// --- Строки, массивы и числа (Длина и размер) ---
			case "min":
				// Подходит и для длины строки, и для значения числа
				msg = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
			case "max":
				msg = fmt.Sprintf("%s must be a maximum of %s", err.Field(), err.Param())
			case "len":
				msg = fmt.Sprintf("%s must have an exact length/value of %s", err.Field(), err.Param())

			// --- Сравнения чисел и времени ---
			case "eq":
				msg = fmt.Sprintf("%s must be equal to %s", err.Field(), err.Param())
			case "ne":
				msg = fmt.Sprintf("%s must not be equal to %s", err.Field(), err.Param())
			case "gt":
				msg = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
			case "gte":
				msg = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
			case "lt":
				msg = fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
			case "lte":
				msg = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())

			// --- Форматы данных и сеть ---
			case "email":
				msg = fmt.Sprintf("%s must be a valid email address", err.Field())
			case "url":
				msg = fmt.Sprintf("%s must be a valid URL", err.Field())
			case "uuid", "uuid3", "uuid4", "uuid5":
				msg = fmt.Sprintf("%s must be a valid UUID", err.Field())
			case "ip", "ipv4", "ipv6":
				msg = fmt.Sprintf("%s must be a valid IP address", err.Field())
			case "mac":
				msg = fmt.Sprintf("%s must be a valid MAC address", err.Field())

			// --- Содержимое строк ---
			case "alpha":
				msg = fmt.Sprintf("%s can only contain alphabetic characters", err.Field())
			case "alphanum":
				msg = fmt.Sprintf("%s can only contain alphanumeric characters", err.Field())
			case "numeric":
				msg = fmt.Sprintf("%s must be a valid numeric value", err.Field())
			case "number":
				msg = fmt.Sprintf("%s must be a valid number", err.Field())
			case "hexadecimal":
				msg = fmt.Sprintf("%s must be a valid hexadecimal string", err.Field())
			case "lowercase":
				msg = fmt.Sprintf("%s must contain only lowercase characters", err.Field())
			case "uppercase":
				msg = fmt.Sprintf("%s must contain only uppercase characters", err.Field())

			// --- Специфичные теги ---
			case "oneof":
				// err.Param() для oneof содержит значения через пробел (например "red green blue")
				// Для красоты вывода меняем пробелы на запятые
				options := strings.ReplaceAll(err.Param(), " ", ", ")
				msg = fmt.Sprintf("%s must be one of: [%s]", err.Field(), options)
			case "datetime":
				msg = fmt.Sprintf("%s must be a valid datetime in the format %s", err.Field(), err.Param())
			case "unique":
				msg = fmt.Sprintf("%s must contain unique values", err.Field())

			// --- Кастомный тег ---
			// case "strong_password":
			// 	msg = fmt.Sprintf("%s must contain at least one special character and number", err.Field())

			// --- Fallback (Если тег не описан выше) ---
			default:
				msg = fmt.Sprintf("%s failed validation on the '%s' tag", err.Field(), err.Tag())
			}

			errors = append(errors, ValidationErrorResponse{
				Field:   err.Field(), // Теперь тут точно JSON-ключ
				Message: msg,
			})
		}
	}

	return errors
}
