package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
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
	json := jsoniter.ConfigFastest
	var user User

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			continue
		}

		email := strings.ToLower(user.Email)
		if !strings.HasSuffix(email, domain) {
			continue
		}

		domainName := strings.SplitN(email, "@", 2)[1]
		result[domainName]++
	}

	return result, scanner.Err()
}
