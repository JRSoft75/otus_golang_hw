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
		m sync.Map // Используем sync.Map для безопасного параллельного доступа
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
	return &DomainStatSync{}
}

func getUsers(r io.Reader) (result users, err error) {
	//var result users
	reader := bufio.NewReader(r)
	for i := 0; i < len(result); i++ {
		line, err := reader.ReadString('\n') // Читаем строку
		if errors.Is(err, io.EOF) {
			break // Достигнут конец файла
		}
		if err != nil {
			return result, fmt.Errorf("error reading line: %w", err)
		}

		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			continue
		}
		result[i] = user
		//result = append(result, user)
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	var wg sync.WaitGroup
	domainStatSync := NewDomainStatSync()

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

				count, _ := domainStatSync.m.LoadOrStore(domainPart, 0) // Получаем текущее значение или создаем новое
				//count, _ := domainStatSync.m.Load(domainPart)
				domainStatSync.m.Store(domainPart, count.(int)+1) // Увеличиваем счетчик
			}
		}(user)
	}
	wg.Wait()
	result := make(map[string]int)
	domainStatSync.m.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(int)
		return true // Продолжаем итерацию
	})
	return result, nil
}
