package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
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

type (
	DomainStat     map[string]int
	DomainStatSync struct {
		mu    sync.RWMutex   // Mutex для безопасного доступа к map
		stats map[string]int // Храним статистику доменов
	}
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

//type users [100_000]User

func NewDomainStatSync() *DomainStatSync {
	return &DomainStatSync{
		stats: make(map[string]int),
	}
}

func getUsers(r io.Reader) ([]User, error) {
	var users []User
	decoder := json.NewDecoder(r)
	for {
		var user User
		err := decoder.Decode(&user)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return users, fmt.Errorf("error decoding user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func countDomains(u []User, domain string) (DomainStat, error) {
	var wg sync.WaitGroup
	result := NewDomainStatSync()
	targetDomain := strings.ToLower(domain)

	for _, user := range u {
		wg.Add(1)
		go func(user User) {
			defer wg.Done()
			email := strings.ToLower(user.Email)
			if strings.Contains(email, targetDomain) {
				emailParts := strings.SplitN(user.Email, "@", 2)
				if len(emailParts) < 2 {
					return // Пропускаем некорректные email
				}
				domainPart := strings.ToLower(emailParts[1])

				result.mu.Lock()
				defer result.mu.Unlock()
				result.stats[domainPart]++
			}
		}(user)
	}
	wg.Wait()
	return result.stats, nil
}
