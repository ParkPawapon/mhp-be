package utils

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func init() {
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{9,15}$`)

	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return field.Name
		}
		return name
	})

	_ = Validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}
		return phoneRegex.MatchString(value)
	})
}

func ValidationErrors(err error) map[string]string {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	out := make(map[string]string)
	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		switch fe.Tag() {
		case "required":
			out[field] = "required"
		case "min":
			out[field] = "minimum"
		case "max":
			out[field] = "maximum"
		case "email":
			out[field] = "invalid email"
		case "phone":
			out[field] = "invalid phone"
		default:
			out[field] = "invalid"
		}
	}
	return out
}
