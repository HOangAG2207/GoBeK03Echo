package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Lấy tên field theo json tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate struct
func ValidateStruct(s any) error {
	return validate.Struct(s)
}

// Format lỗi đẹp dạng list
func FormatValidationErrors(err error) []FieldError {
	var errorsList []FieldError

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errorsList
	}

	for _, e := range validationErrors {
		field := e.Field()
		tag := e.Tag()
		param := e.Param()

		// 1. Lấy message template
		msgTemplate, ok := Messages[tag]
		if !ok {
			msgTemplate = "không hợp lệ"
		}

		// 2. Replace param nếu có
		msg := strings.ReplaceAll(msgTemplate, "{param}", param)

		errorsList = append(errorsList, FieldError{
			Field:   field,
			Message: msg,
		})
	}

	return errorsList
}
