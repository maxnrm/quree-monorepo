package sendlimiter

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

type UserRateLimiter struct {
	ChatID      string
	RateLimiter *rate.Limiter
	LastMsgSent time.Time
}

type SendLimiter struct {
	Ctx                   context.Context
	GlobalRateLimiter     *rate.Limiter
	UserRateLimitersCache map[string]*UserRateLimiter
}

func Init(ctx context.Context, gl int, gb int) *SendLimiter {
	limit := rate.Every(time.Second / time.Duration(gl))
	rateLimiter := rate.NewLimiter(limit, gb)

	return &SendLimiter{
		Ctx:                   ctx,
		GlobalRateLimiter:     rateLimiter,
		UserRateLimitersCache: make(map[string]*UserRateLimiter),
	}

}

func (sl *SendLimiter) AddUserRateLimiter(chatID string, l int, b int) {
	limit := rate.Every(time.Second / time.Duration(l))
	rateLimiter := rate.NewLimiter(limit, b)

	sl.UserRateLimitersCache[chatID] = &UserRateLimiter{
		ChatID:      chatID,
		RateLimiter: rateLimiter,
		LastMsgSent: time.Now(),
	}
}

func (sl *SendLimiter) GetUserRateLimiter(chatID string) *UserRateLimiter {
	if v, ok := sl.UserRateLimitersCache[chatID]; ok {
		return v
	}

	return nil
}

func (sl *SendLimiter) removeUserRateLimiter(chatID string) {
	delete(sl.UserRateLimitersCache, chatID)
}

func (sl *SendLimiter) RemoveOldUserRateLimitersCache(delay time.Duration) {
	for {
		time.Sleep(delay * time.Second)
		for k, v := range sl.UserRateLimitersCache {
			if time.Since(v.LastMsgSent) > delay*time.Second {
				fmt.Println(v.LastMsgSent, "removing...")
				fmt.Println(time.Now())
				sl.removeUserRateLimiter(k)
			}
		}
		// fmt.Println("Clearing cache, items:", sl.UserRateLimitersCache)
		// fmt.Println("Num goroutine:", runtime.NumGoroutine())
	}
}
