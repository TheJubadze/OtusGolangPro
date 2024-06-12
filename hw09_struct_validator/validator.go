package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result strings.Builder
	for _, err := range v {
		result.WriteString(fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error()))
	}
	return result.String()
}

func Validate(val interface{}) error {
	v := reflect.ValueOf(val)
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("val must be a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath == "" {
			tag := field.Tag.Get("validate")
			if tag != "" {
				err := validateField(field.Name, v.Field(i).Interface(), tag)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(fieldName string, field any, tag string) error {
	switch v := reflect.ValueOf(field); v.Kind() {
	case reflect.Int, reflect.String:
		return validateValue(fieldName, v, tag)
	case reflect.Slice:
		return validateSlice(fieldName, v, tag)
	default:
		return fmt.Errorf("unsupported type for field '%s': %T", fieldName, field)
	}
}
