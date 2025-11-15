package middleware

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(data interface{}) error {
	validate := validator.New()
	err := validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		str, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return strings.TrimSpace(str) != ""
	})
	if err != nil {
		return err
	}
	err = validate.Struct(data)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		if len(validationErrors) > 0 {
			fieldErr := validationErrors[0]
			fieldName := strings.ToLower(fieldErr.Field())

			switch fieldErr.Tag() {
			case "required":
				err = fmt.Errorf("%s is required", fieldName)
			case "min":
				if fieldErr.Kind() == reflect.String {
					err = fmt.Errorf("%s must be at least %s characters", fieldName, fieldErr.Param())
				} else {
					err = fmt.Errorf("%s must be at least %s", fieldName, fieldErr.Param())
				}
			case "max":
				if fieldErr.Kind() == reflect.String {
					err = fmt.Errorf("%s must be at most %s characters", fieldName, fieldErr.Param())
				} else {
					err = fmt.Errorf("%s must be at most %s", fieldName, fieldErr.Param())
				}
			case "uuid4":
				err = fmt.Errorf("%s must be a valid UUID", fieldName)
			default:
				err = fmt.Errorf("%s is invalid", fieldName)
			}

			return err
		}
	}

	return nil
}
