package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-framework-v2/go-access/access"
	"go.uber.org/zap"
)

// IPLimit 按 IP 的滑动窗口限流：每个 key 在 window 内最多 max 次请求
type IPLimit struct {
	mu      sync.Mutex
	records map[string][]time.Time // key -> 最近请求时间
	max     int                    // 窗口内允许的最大次数
	window  time.Duration          // 时间窗口，如 1 * time.Minute
}

// NewIPLimit max=30, window=1*time.Minute 表示每 IP 每分钟最多 30 次
func NewIPLimit(max int, window time.Duration) *IPLimit {
	return &IPLimit{
		records: make(map[string][]time.Time),
		max:     max,
		window:  window,
	}
}

// Key 从 c 取限流 key，默认用 ClientIP
func (l *IPLimit) Key(c *gin.Context) string {
	return c.ClientIP()
}

// Allow 返回是否允许本次请求；若不允许会写 429 并 abort
func (l *IPLimit) Allow(c *gin.Context) bool {
	key := l.Key(c)
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	// 只保留在 window 内的请求时间
	cut := now.Add(-l.window)
	var kept []time.Time
	for _, t := range l.records[key] {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= l.max {
		l.records[key] = kept
		return false
	}
	l.records[key] = append(kept, now)
	return true
}

// Middleware 返回 gin 中间件：超出限制时 429 并 abort
func (l *IPLimit) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		emptyInterfaceData := interface{}(nil)
		if !l.Allow(c) {
			zap.S().Warn("rate limit exceeded", zap.String("ip", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusOK, access.GetErrorResult(http.StatusTooManyRequests, emptyInterfaceData, "请求过于频繁，请稍后再试"))
			return
		}
		c.Next()
	}
}
