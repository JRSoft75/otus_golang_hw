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
				email := strings.ToLower(user.Email)
				if strings.HasSuffix(email, targetDomain) {
					if parts := strings.SplitN(email, "@", 2); len(parts) == 2 {
						domainStatSync.Increment(parts[1])
					}
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
