package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)
import jsoniter "github.com/json-iterator/go"

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
	for scanner.Scan() {
		email := strings.ToLower(jsoniter.Get(scanner.Bytes(), "Email").ToString())
		if email == "" || !strings.HasSuffix(email, "."+domain) {
			continue
		}

		parts := strings.SplitN(email, "@", 2)
		if len(parts) != 2 {
			continue
		}

		result[parts[1]]++
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return result, nil
}
