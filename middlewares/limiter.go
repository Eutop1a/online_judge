package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimiterMiddleWare 限流中间件, Token Bucket 算法来创建了一个令牌桶
// fillInterval：填充令牌的时间间隔，即每隔多久会向桶中放入一定数量的令牌，这里是1秒，表示每秒会向桶中放入一定数量的令牌。
// cap：令牌桶的容量，即桶中最多能存储多少令牌。
// quantum：每次填充的令牌数量，即每次放入多少个令牌。
func RateLimiterMiddleWare(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum) // bucket.TakeAvailable(1) 来尝试获取1个令牌，如果成功获取到令牌（即返回的值大于等于1）
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
