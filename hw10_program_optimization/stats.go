//go:generate easyjson -all stats.go

package hw10programoptimization

import (
	"bufio"
	"encoding/json"
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

func GetDomainStatOld(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsersOpt(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomainsOpt(u, domain), nil
}

type usersOpt []*User

func getUsersOpt(r io.Reader) (usersOpt, error) {
	var result usersOpt
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		user := new(User)
		line := scanner.Bytes()
		if err := user.UnmarshalJSON(line); err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}

func countDomainsOpt(u usersOpt, domain string) DomainStat {
	result := make(DomainStat)
	domain = "." + domain

	for _, user := range u {
		if !strings.HasSuffix(user.Email, domain) {
			continue
		}

		at := strings.LastIndexByte(user.Email, '@')
		if at < 0 {
			continue
		}

		result[strings.ToLower(user.Email[at+1:])]++
	}
	return result
}
