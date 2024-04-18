package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"
	"golang.org/x/example/hello/reverse"
)

const shield = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if err := validateString(str); err != nil {
		return "", err
	}

	return unpack(str)
}

func validateString(str string) error {
	for i := 0; i < len(str)-1; i++ {
		if rune(str[i]) == shield {
			i++
			continue
		}
		_, err0 := strconv.Atoi(string(str[i]))
		_, err1 := strconv.Atoi(string(str[i+1]))
		if err0 == nil && err1 == nil {
			return ErrInvalidString
		}
	}
	return nil
}

func unpack(str string) (string, error) {
	st := stack.New()
	backslash := false

	for i, ch := range str {
		char := string(ch)
		if backslash {
			if _, err := strconv.Atoi(char); err == nil || ch == shield {
				st.Push(ch)
				backslash = false
			} else {
				return "", ErrInvalidString
			}
			continue
		}

		if ch == shield {
			if i == len(str)-1 {
				return "", ErrInvalidString
			}
			backslash = true
			continue
		}

		if n, err := strconv.Atoi(char); err == nil {
			if st.Len() < 1 {
				return "", ErrInvalidString
			}

			if n < 1 {
				st.Pop()
			} else {
				repCh := st.Peek().(rune)
				for i := 0; i < n-1; i++ {
					st.Push(repCh)
				}
			}
			continue
		}

		st.Push(ch)
	}

	sb := strings.Builder{}

	for st.Len() > 0 {
		sb.WriteRune(st.Pop().(rune))
	}

	return reverse.String(sb.String()), nil
}
