package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

func validateValue(fieldName string, value reflect.Value, tag string) []Error {
	validationFunc, err := getValidationFunc(fieldName, value)
	if err != nil {
		return []Error{&ProgramError{Err: err}}
	}

	constraints := strings.Split(tag, "|")
	var errors []Error
	for _, constraint := range constraints {
		parts := strings.SplitN(constraint, ":", 2)
		if len(parts) != 2 {
			errors = append(errors, &TagError{Value: fieldName, Err: fmt.Errorf("invalid tag format for field '%s'", fieldName)})
			continue
		}
		err := validationFunc(value, parts)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func getValidationFunc(fieldName string, value reflect.Value) (func(value reflect.Value, parts []string) Error, error) {
	var validationFunc func(value reflect.Value, parts []string) Error
	switch value.Kind() {
	case reflect.Int:
		validationFunc = func(value reflect.Value, parts []string) Error {
			return validateInt(fieldName, int(value.Int()), parts)
		}
	case reflect.String:
		validationFunc = func(value reflect.Value, parts []string) Error {
			return validateString(fieldName, value.String(), parts)
		}
	default:
		return nil, fmt.Errorf("unsupported type for field '%s': %T", fieldName, value)
	}
	return validationFunc, nil
}
