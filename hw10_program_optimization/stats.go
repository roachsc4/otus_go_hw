package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
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

var emailRegexp = regexp.MustCompile(`.+@\w+(\.\w+)+$`)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain), nil
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)

	for i := 0; scanner.Scan(); i++ {
		var user User
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) DomainStat {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, "."+domain) && emailRegexp.MatchString(user.Email) {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result
}
