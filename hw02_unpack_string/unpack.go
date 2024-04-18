package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"
	"golang.org/x/example/hello/reverse"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	for i := 1; i < len(str); i++ {
		if string(str[i-1]) == `\` {
			i += 2
			continue
		}
		_, err0 := strconv.Atoi(string(str[i-1]))
		_, err1 := strconv.Atoi(string(str[i]))
		if err0 == nil && err1 == nil {
			return "", ErrInvalidString
		}
	}

	return processString(str)
}

func processString(str string) (string, error) {
	st := stack.New()
	backslash := false

	for _, ch := range str {
		if backslash {
			if _, err := strconv.Atoi(string(ch)); err == nil || string(ch) == `\` {
				st.Push(string(ch))
				backslash = false
			} else {
				return "", ErrInvalidString
			}
			continue
		}

		if string(ch) == `\` {
			backslash = true
			continue
		}

		if n, err := strconv.Atoi(string(ch)); err == nil {
			if st.Len() < 1 {
				return "", ErrInvalidString
			}

			if n < 1 {
				st.Pop()
			} else {
				repCh := st.Peek().(string)
				for i := 0; i < n-1; i++ {
					st.Push(repCh)
				}
			}
			continue
		}

		st.Push(string(ch))
	}

	sb := strings.Builder{}

	for st.Len() > 0 {
		sb.WriteString(st.Pop().(string))
	}

	return reverse.String(sb.String()), nil
}
