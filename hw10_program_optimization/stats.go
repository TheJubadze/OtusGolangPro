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
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		if err = json.Unmarshal(scanner.Bytes(), &result[i]); err != nil {
			return
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = "." + domain

	for _, user := range u {
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
