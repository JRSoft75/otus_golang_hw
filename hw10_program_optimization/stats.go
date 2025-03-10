package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
)

type User struct {
	Email string `json:"email"`
}

type (
	DomainStat     map[string]int
	DomainStatSync struct {
		sync.Mutex
		stats map[string]int
	}
)

// Используем пул объектов для повторного использования структуры User.
var userPool = sync.Pool{
	New: func() interface{} { return new(User) },
}

var domainPool = sync.Pool{
	New: func() interface{} { return new(string) },
}

func (d *DomainStatSync) Increment(domain string) {
	d.Lock()
	defer d.Unlock()
	d.stats[domain]++
}

func (d *DomainStatSync) ToMap() DomainStat {
	d.Lock()
	defer d.Unlock()
	result := make(DomainStat, len(d.stats))
	for k, v := range d.stats {
		result[k] = v
	}
	return result
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var wg sync.WaitGroup
	jobs := make(chan User, 1024)
	domainStatSync := &DomainStatSync{stats: make(map[string]int)}
	targetDomain := strings.ToLower(domain)
	workerCount := runtime.NumCPU()
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range jobs {
				email := user.Email
				atIndex := strings.LastIndex(email, "@")
				if atIndex == -1 || atIndex == len(email)-1 {
					continue
				}

				domainPart := email[atIndex+1:]
				domainPartLower := strings.ToLower(domainPart)

				if len(domainPartLower) < len(targetDomain) {
					continue
				}

				if domainPartLower[len(domainPartLower)-len(targetDomain):] == targetDomain {
					domain := domainPool.Get().(*string)
					*domain = domainPartLower
					domainStatSync.Increment(*domain)
					domainPool.Put(domain)
				}
			}
		}()
	}

	decoder := json.NewDecoder(r)
	var user *User
	for {
		user = userPool.Get().(*User)
		err := decoder.Decode(user)
		if err != nil {
			userPool.Put(user)
			close(jobs)
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decoding error: %w", err)
		}
		jobs <- *user
		userPool.Put(user)
	}

	wg.Wait()
	return domainStatSync.ToMap(), nil
}
