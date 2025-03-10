package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
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
		stats sync.Map
	}
)

func (d *DomainStatSync) Increment(domain string) {
	val, _ := d.stats.LoadOrStore(domain, &atomic.Int64{})
	counter := val.(*atomic.Int64)
	counter.Add(1)
}

func (d *DomainStatSync) ToMap() DomainStat {
	result := make(DomainStat)
	d.stats.Range(func(key, value any) bool {
		result[key.(string)] = int(value.(*atomic.Int64).Load())
		return true
	})
	return result
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	users, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(users, domain)
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
	domainStatSync := &DomainStatSync{}
	targetDomain := strings.ToLower(domain)

	workerCount := 4 // Количество воркеров
	jobs := make(chan User, workerCount*100)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range jobs {
				email := strings.ToLower(user.Email)
				if strings.HasSuffix(email, targetDomain) {
					emailParts := strings.SplitN(email, "@", 2)
					if len(emailParts) == 2 {
						domainPart := emailParts[1]
						domainStatSync.Increment(domainPart)
					}
				}
			}
		}()
	}

	for _, user := range u {
		jobs <- user
	}
	close(jobs)

	wg.Wait()
	return domainStatSync.ToMap(), nil
}
