package hw09structvalidator

import (
	"fmt"
	"reflect"
)

func validateSlice(fieldName string, slice reflect.Value, tag string) []Error {
	var errors []Error
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i).Interface()
		if err := validateField(fmt.Sprintf("%s[%d]", fieldName, i), elem, tag); err != nil {
			errors = append(errors, err...)
		}
	}

	return errors
}
