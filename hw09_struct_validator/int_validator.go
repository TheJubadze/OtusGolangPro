package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

func validateInt(i int, parts []string) error {
	key := parts[0]
	valueStr := parts[1]

	switch key {
	case "min":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", valueStr)
		}
		if i < value {
			return fmt.Errorf("value %d is less than min %d", i, value)
		}
	case "max":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", valueStr)
		}
		if i > value {
			return fmt.Errorf("value %d is greater than max %d", i, value)
		}
	case "in":
		inValues := strings.Split(valueStr, ",")
		for _, str := range inValues {
			value, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("invalid in value: %s", str)
			}
			if i == value {
				return nil
			}
		}
		return fmt.Errorf("value %d is not in %s", i, valueStr)
	default:
		return fmt.Errorf("unknown validation key: %s", key)
	}
	return nil
}
