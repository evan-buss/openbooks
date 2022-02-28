package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors      map[string]*visitor
	visitorsMutex sync.Mutex
	timeInterval  time.Duration
	burstInterval int
}

func NewRateLimiter(timeInterval time.Duration, burstInterval int) *RateLimiter {
	return &RateLimiter{
		visitors:      make(map[string]*visitor),
		timeInterval:  timeInterval,
		burstInterval: burstInterval,
	}
}

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.visitorsMutex.Lock()
	defer rl.visitorsMutex.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(rl.timeInterval), rl.burstInterval)
		// Include the current time when creating a new visitor.
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func (rl *RateLimiter) CleanupVisitors() {
	rl.visitorsMutex.Lock()
	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 1*time.Minute {
			log.Printf("Removing %s from rate limit cache.\n", ip)
			delete(rl.visitors, ip)
		}
	}
	rl.visitorsMutex.Unlock()
}

func (server *server) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Read forwarded headers as well...
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := server.rateLimiter.getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
