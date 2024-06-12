package hw09structvalidator

import (
	"fmt"
	"reflect"
)

func validateSlice(fieldName string, slice reflect.Value, tag string) error {
	var validationErrors []error
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i).Interface()
		if err := validateField(fmt.Sprintf("%s[%d]", fieldName, i), elem, tag); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("index %d: %w", i, err))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("slice validation errors: %v", validationErrors)
	}
	return nil
}
