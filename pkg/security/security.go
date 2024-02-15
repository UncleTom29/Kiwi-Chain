package main 

type RateLimiter struct {
	visitors map[string]*time.Timer
	mu       sync.Mutex
	Rate     time.Duration
}

func NewRateLimiter(rate time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*time.Timer),
		Rate:     rate,
	}
}

func (rl *RateLimiter) Limit(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, exists := rl.visitors[ip]; !exists {
		rl.visitors[ip] = time.AfterFunc(rl.Rate, func() {
			rl.mu.Lock()
			defer rl.mu.Unlock()
			delete(rl.visitors, ip)
		})
		return true
	}

	return false
}
