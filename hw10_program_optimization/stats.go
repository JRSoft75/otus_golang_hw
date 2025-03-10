package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
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

// Используем пул объектов для повторного использования структуры User.
var userPool = sync.Pool{
	New: func() interface{} { return new(User) },
}

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
	var wg sync.WaitGroup
	jobs := make(chan User, 100)
	domainStatSync := &DomainStatSync{}
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
