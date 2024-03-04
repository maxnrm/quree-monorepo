package sendlimiter

import (
	"context"
	"fmt"
	"quree/config"
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

func Init(ctx context.Context) *SendLimiter {
	limit := rate.Every(time.Second / time.Duration(config.RATE_LIMIT_GLOBAL))
	rateLimiter := rate.NewLimiter(limit, config.RATE_LIMIT_BURST_GLOBAL)

	return &SendLimiter{
		Ctx:                   ctx,
		GlobalRateLimiter:     rateLimiter,
		UserRateLimitersCache: make(map[string]*UserRateLimiter),
	}

}

func (sl *SendLimiter) AddUserRateLimiter(chatID string) {
	limit := rate.Every(time.Second / time.Duration(config.RATE_LIMIT_USER))
	rateLimiter := rate.NewLimiter(limit, config.RATE_LIMIT_BURST_USER)

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

func (sl *SendLimiter) RemoveOldUserRateLimitersCache() {
	for {
		time.Sleep(10 * time.Second)
		for k, v := range sl.UserRateLimitersCache {
			if time.Since(v.LastMsgSent) > 10*time.Second {
				fmt.Println(v.LastMsgSent, "removing...")
				fmt.Println(time.Now())
				sl.removeUserRateLimiter(k)
			}
		}
		// fmt.Println("Clearing cache, items:", sl.UserRateLimitersCache)
		// fmt.Println("Num goroutine:", runtime.NumGoroutine())
	}
}
