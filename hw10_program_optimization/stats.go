package hw10programoptimization

import (
	"bufio"
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

type users [100_000]User

func NewDomainStatSync() *DomainStatSync {
	return &DomainStatSync{
		stats: make(map[string]int),
	}
}

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewReader(r)
	for i := 0; i < len(result); i++ {
		line, err := reader.ReadString('\n') // Читаем строку
		if errors.Is(err, io.EOF) {
			break // Достигнут конец файла
		}
		if err != nil {
			return result, err
		}

		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return result, err
		}
		result[i] = user
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	var wg sync.WaitGroup
	result := NewDomainStatSync()

	for _, user := range u {
		wg.Add(1)
		go func(user User) {
			defer wg.Done()
			if strings.Contains(user.Email, domain) {
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
