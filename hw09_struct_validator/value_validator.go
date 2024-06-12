package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

func validateValue(fieldName string, value reflect.Value, tag string) error {
	validationFunc, err := getValidationFunc(fieldName, value)
	if err != nil {
		return err
	}

	constraints := strings.Split(tag, "|")
	var validationErrors []error
	for _, constraint := range constraints {
		parts := strings.SplitN(constraint, ":", 2)
		if len(parts) != 2 {
			validationErrors = append(validationErrors, fmt.Errorf("invalid tag format for field '%s'", fieldName))
			continue
		}
		err := validationFunc(value, parts)
		if err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors: %v", validationErrors)
	}
	return nil
}

func getValidationFunc(fieldName string, value reflect.Value) (func(value reflect.Value, parts []string) error, error) {
	var validationFunc func(value reflect.Value, parts []string) error
	switch value.Kind() {
	case reflect.Int:
		validationFunc = func(value reflect.Value, parts []string) error {
			return validateInt(int(value.Int()), parts)
		}
	case reflect.String:
		validationFunc = func(value reflect.Value, parts []string) error {
			return validateString(value.String(), parts)
		}
	default:
		return nil, fmt.Errorf("unsupported type for field '%s': %T", fieldName, value)
	}
	return validationFunc, nil
}
