package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"

	"football-api/internal/core/config"
	"football-api/internal/shared/response"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = map[string]*visitor{}
	mu       sync.Mutex
	redisMu  sync.Mutex
	rdb      *redis.Client
)

func getVisitor(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(r, b)
		visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func redisClient(cfg config.Config) *redis.Client {
	redisMu.Lock()
	defer redisMu.Unlock()
	if rdb != nil {
		return rdb
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return rdb
}

func LoginRateLimit(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.RateLimitStore == "redis" {
			ctx := context.Background()
			key := fmt.Sprintf("rl:login:%s", c.ClientIP())
			window := time.Duration(cfg.LoginWindowMin) * time.Minute

			count, err := redisClient(cfg).Incr(ctx, key).Result()
			if err == nil {
				if count == 1 {
					_ = redisClient(cfg).Expire(ctx, key, window).Err()
				}
				if int(count) > cfg.LoginLimit {
					response.Error(c, http.StatusTooManyRequests, "too many login attempts", "RATE_LIMITED")
					c.Abort()
					return
				}
				c.Next()
				return
			}
		}

		interval := time.Duration(cfg.LoginWindowMin) * time.Minute
		limiter := getVisitor(c.ClientIP(), rate.Every(interval), cfg.LoginLimit)
		if !limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, "too many login attempts", "RATE_LIMITED")
			c.Abort()
			return
		}
		c.Next()
	}
}
