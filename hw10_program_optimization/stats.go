package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewReader(r)
	for i := 0; i < len(result); i++ {
		line, err := reader.ReadString('\n') // Читаем строку
		if err == io.EOF {
			break // Достигнут конец файла
		}
		if err != nil {
			return result, nil
		}

		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return result, nil
		}
		result[i] = user
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, domain) {
			emailParts := strings.SplitN(user.Email, "@", 2)
			if len(emailParts) < 2 {
				continue // Пропускаем некорректные email
			}
			result[strings.ToLower(emailParts[1])]++
		}
	}
	return result, nil
}
