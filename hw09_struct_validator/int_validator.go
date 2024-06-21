package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

func validateInt(fieldName string, i int, parts []string) Error {
	key := parts[0]
	valueStr := parts[1]

	switch key {
	case "min":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return &TagError{Key: key, Value: valueStr, Err: err}
		}
		if i < value {
			return &ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("value %d is less than min %d", i, value),
			}
		}
	case "max":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return &TagError{Key: key, Value: valueStr, Err: err}
		}
		if i > value {
			return &ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("value %d is greater than max %d", i, value),
			}
		}
	case "in":
		inValues := strings.Split(valueStr, ",")
		for _, str := range inValues {
			value, err := strconv.Atoi(str)
			if err != nil {
				return &TagError{Key: key, Value: str, Err: err}
			}
			if i == value {
				return nil
			}
		}
		return &ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("value %d is not in %s", i, valueStr),
		}
	default:
		return &TagError{Key: key, Value: "unknown key", Err: fmt.Errorf("unknown key")}
	}
	return nil
}
