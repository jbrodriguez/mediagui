package lib

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	tokens     chan struct{}
	ticker     *time.Ticker
	mu         sync.RWMutex
	capacity   int
	refillRate time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewRateLimiter creates a new rate limiter with the specified capacity and refill rate
// capacity: maximum number of tokens in the bucket
// refillRate: how often to add a token to the bucket
func NewRateLimiter(capacity int, refillRate time.Duration) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())

	rl := &RateLimiter{
		tokens:     make(chan struct{}, capacity),
		ticker:     time.NewTicker(refillRate),
		capacity:   capacity,
		refillRate: refillRate,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Fill the bucket initially
	for range capacity {
		select {
		case rl.tokens <- struct{}{}:
		default:
		}
	}

	// Start the refill goroutine
	go rl.refill()

	return rl
}

// Wait blocks until a token is available, then consumes it
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// TryWait attempts to consume a token without blocking
// Returns true if successful, false if no token is available
func (rl *RateLimiter) TryWait() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// WaitWithTimeout waits for a token with a timeout
// Returns true if a token was acquired, false if timeout occurred
func (rl *RateLimiter) WaitWithTimeout(timeout time.Duration) bool {
	select {
	case <-rl.tokens:
		return true
	case <-time.After(timeout):
		return false
	}
}

// refill adds tokens to the bucket at the specified rate
func (rl *RateLimiter) refill() {
	defer rl.ticker.Stop()

	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-rl.ticker.C:
			// Try to add a token to the bucket
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Bucket is full, skip
			}
		}
	}
}

// Stop stops the rate limiter and cleans up resources
func (rl *RateLimiter) Stop() {
	rl.cancel()
}

// Available returns the number of available tokens (approximate)
func (rl *RateLimiter) Available() int {
	return len(rl.tokens)
}
