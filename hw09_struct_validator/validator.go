package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

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
				errs := validateField(field.Name, v.Field(i).Interface(), tag)
				for _, err := range errs {
					var e *ValidationError
					if errors.As(err, &e) && e.IsValidation() {
						validationErrors = append(validationErrors, *e)
					} else {
						log.Default().Printf("error validating field '%s': %v", field.Name, err)
					}
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(fieldName string, field any, tag string) []Error {
	switch v := reflect.ValueOf(field); v.Kind() {
	case reflect.Int, reflect.String:
		return validateValue(fieldName, v, tag)
	case reflect.Slice:
		return validateSlice(fieldName, v, tag)
	default:
		return []Error{&TagError{Key: fieldName, Value: tag, Err: fmt.Errorf("unsupported type for field '%s': %T", fieldName, field)}}
	}
}
