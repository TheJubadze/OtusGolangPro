package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func validateString(fieldName string, s string, parts []string) Error {
	key := parts[0]
	valueStr := parts[1]

	switch key {
	case "len":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return &TagError{Key: key, Value: valueStr, Err: err}
		}
		if len(s) != value {
			return &ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("length of '%s' is %d, expected %d", s, len(s), value),
			}
		}
	case "regexp":
		re, err := regexp.Compile(valueStr)
		if err != nil {
			return &TagError{Key: key, Value: valueStr, Err: err}
		}
		if !re.MatchString(s) {
			return &ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("'%s' does not match pattern %s", s, valueStr),
			}
		}
	case "in":
		inValues := strings.Split(valueStr, ",")
		for _, value := range inValues {
			if s == value {
				return nil
			}
		}
		return &ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("'%s' is not in %s", s, valueStr),
		}
	default:
		return &TagError{Key: key, Value: "unknown key", Err: fmt.Errorf("unknown key")}
	}
	return nil
}
