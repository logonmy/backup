package server

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type rateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter .
func NewIPRateLimiter(r rate.Limit, b int) *rateLimiter {
	i := &rateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
	go func() {
		for {
			select {
			case <-time.After(time.Minute * 5):
				i.ips = make(map[string]*rate.Limiter)
			}
		}
	}()
	return i
}

// addIP 创建了一个新的速率限制器，并将其添加到 ips 映射中,
// 使用 IP地址作为密钥
func (i *rateLimiter) addIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limiter

	return limiter
}

// GetLimiter 返回所提供的IP地址的速率限制器(如果存在的话).
// 否则调用 addIP 将 IP 地址添加到映射中
func (i *rateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.addIP(ip)
	}

	i.mu.Unlock()

	return limiter
}

