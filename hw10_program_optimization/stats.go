package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		var user User
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return nil, err
		}
		if strings.Contains(user.Email, "."+domain) && strings.Contains(user.Email, "@") {
			splittedString := strings.SplitN(user.Email, "@", 2)
			if len(splittedString) == 2 {
				result[strings.ToLower(splittedString[1])]++
			}
		}
	}

	return result, nil
}
