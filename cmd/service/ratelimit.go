package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"

	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	cache     map[string]*visitor
	cacheLock *sync.Mutex
}

type visitor struct {
	limiter    ratelimit.Limiter
	lastAccess time.Time
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{cache: make(map[string]*visitor), cacheLock: &sync.Mutex{}}
	go rl.cleanup(time.Second * 30)
	return rl
}

func newVisitor() *visitor {
	return &visitor{
		limiter:    ratelimit.New(1),
		lastAccess: time.Now(),
	}
}

func (rl *RateLimiter) Access(address string) {
	rl.cacheLock.Lock()
	log.Println(address)
	val, exists := rl.cache[address]
	if !exists {
		log.Println("Adding new visitor?")
		val = newVisitor()
		rl.cache[address] = val
	}
	val.lastAccess = time.Now()
	rl.cacheLock.Unlock()

	val.limiter.Take()
}

func (rl *RateLimiter) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		ipAddr, err := getIPAddressFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unavailable, "Unable to determine IP address.")
		}

		rl.Access(ipAddr)

		return handler(ctx, req)
	}
}

func getIPAddressFromContext(ctx context.Context) (string, error) {
	// If request is HTTP use the x-forwarded-for
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md["x-forwarded-for"]) != 0 {
		return md["x-forwarded-for"][0], nil
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		return p.Addr.String(), nil
	}

	return "", errors.New("unable to determine IP address from context")
}

// cleanup runs in a goroutine and removes old entries from the rate limiter cache
func (rl *RateLimiter) cleanup(duration time.Duration) {
	for {
		rl.cacheLock.Lock()

		for ip, visitor := range rl.cache {
			if time.Since(visitor.lastAccess) > 5*time.Minute {
				delete(rl.cache, ip)
			}
		}

		rl.cacheLock.Unlock()
		time.Sleep(duration)
	}
}
