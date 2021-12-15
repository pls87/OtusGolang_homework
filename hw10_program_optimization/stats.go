package hw10programoptimization

import (
	"bufio"
	"encoding/json"
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	return countDomains(u, err, domain)
}

func getUsers(r io.Reader) (userCh chan User, e chan error) {
	userCh = make(chan User)
	e = make(chan error)
	go func() {
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)
		var user User
		for scanner.Scan() {
			if err := json.Unmarshal([]byte(scanner.Text()), &user); err != nil {
				continue
			}
			userCh <- user
		}
		e <- scanner.Err()

		close(e)
		close(userCh)
	}()

	return userCh, e
}

func countDomains(userCh chan User, e chan error, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for {
		select {
		case u := <-userCh:
			matched, err := regexp.Match("\\."+domain, []byte(u.Email))
			if err != nil {
				return nil, err
			}

			if matched {
				num := result[strings.ToLower(strings.SplitN(u.Email, "@", 2)[1])]
				num++
				result[strings.ToLower(strings.SplitN(u.Email, "@", 2)[1])] = num
			}
		case err := <-e:
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
}
