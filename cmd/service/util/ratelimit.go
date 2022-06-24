package util

import (
	"log"
	"sync"
	"time"

	"go.uber.org/ratelimit"
)

type rateLimiter struct {
	cache map[string]*visitor
	lock  *sync.Mutex
}

func NewRateLimiter() *rateLimiter {
	rl := &rateLimiter{cache: make(map[string]*visitor), lock: &sync.Mutex{}}
	go rl.cleanup(time.Second * 30)
	return rl
}

func newVisitor() *visitor {
	return &visitor{
		limiter:    ratelimit.New(1),
		lastAccess: time.Now(),
	}
}

type visitor struct {
	limiter    ratelimit.Limiter
	lastAccess time.Time
}

func (rl *rateLimiter) Access(address string) {
	rl.lock.Lock()
	log.Println(address)
	val, exists := rl.cache[address]
	if !exists {
		log.Println("Adding new visitor?")
		val = newVisitor()
		rl.cache[address] = val
	}

	val.limiter.Take()
	val.lastAccess = time.Now()

	rl.lock.Unlock()
}

func (rl *rateLimiter) cleanup(duration time.Duration) {
	for {
		rl.lock.Lock()

		for ip, visitor := range rl.cache {
			if time.Since(visitor.lastAccess) > 3*time.Minute {
				delete(rl.cache, ip)
			}
		}

		rl.lock.Unlock()
		time.Sleep(duration)
	}
}
