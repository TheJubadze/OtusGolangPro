package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/goccy/go-json"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	user := User{}
	result := make(DomainStat)
	domain = "." + domain

	for i := 0; scanner.Scan(); i++ {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}
		email := user.Email
		atIndex := strings.LastIndex(email, "@")
		if atIndex == -1 {
			continue
		}
		emailDomain := email[atIndex+1:]

		if hasSuffix(emailDomain, domain) {
			lowerEmailDomain := strings.ToLower(emailDomain)
			result[lowerEmailDomain]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return result, nil
}

func hasSuffix(s, suffix string) bool {
	sLen := len(s)
	suffixLen := len(suffix)
	if suffixLen > sLen {
		return false
	}
	for i := 0; i < suffixLen; i++ {
		if unicode.ToLower(rune(s[sLen-suffixLen+i])) != unicode.ToLower(rune(suffix[i])) {
			return false
		}
	}
	return true
}
