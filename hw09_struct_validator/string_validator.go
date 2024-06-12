package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func validateString(s string, parts []string) error {
	key := parts[0]
	valueStr := parts[1]

	switch key {
	case "len":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fmt.Errorf("invalid len value: %s", valueStr)
		}
		if len(s) != value {
			return fmt.Errorf("length of '%s' is %d, expected %d", s, len(s), value)
		}
	case "regexp":
		re, err := regexp.Compile(valueStr)
		if err != nil {
			return fmt.Errorf("invalid regexp pattern: %s", valueStr)
		}
		if !re.MatchString(s) {
			return fmt.Errorf("'%s' does not match pattern %s", s, valueStr)
		}
	case "in":
		inValues := strings.Split(valueStr, ",")
		for _, value := range inValues {
			if s == value {
				return nil
			}
		}
		return fmt.Errorf("'%s' is not in %s", s, valueStr)
	default:
		return fmt.Errorf("unknown validation key: %s", key)
	}
	return nil
}
