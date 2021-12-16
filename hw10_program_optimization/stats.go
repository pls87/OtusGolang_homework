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
		email := extractEmail(scanner.Bytes())
		handleEmail(email, domain, &result)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return result, nil
}

func extractEmail(bytes []byte) string {
	return strings.ToLower(jsoniter.Get(bytes, "Email").ToString())
}

func handleEmail(email, domain string, stat *DomainStat) {
	if !strings.HasSuffix(email, "."+domain) {
		return
	}

	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return
	}

	(*stat)[parts[1]]++
}
