package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/hohice/gin-web/server/ex"

	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
)

var limiter *redis_rate.Limiter

func init() {
	registerSelf(func(conf *setting.Configs) (error, Closeble) {
		ring := redis.NewRing(&redis.RingOptions{
			Addrs: conf.Limit.AddrMap,
		},
		)
		limiter := redis_rate.NewLimiter(ring)
		// Optional.
		limiter.Fallback = rate.NewLimiter(rate.Every(time.Second), conf.Limit.DefaultRate)

		return nil, func() {}
	})
}

//GetLimitfactor type define method to get factor used by limit
type GetLimitfactor func(c *gin.Context, limit int64) (string, int64)

//DefaultLimitfactor default func of GetLimitfactor
func DefaultLimitfactor(c *gin.Context, limit int64) (string, int64) {
	path := c.Request.URL.Path
	return path, limit
}

//Limiter use getLimitfactor to get fctor and limit then to run
func Limiter(getLimitfactor GetLimitfactor, limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		factor, limit := getLimitfactor(c, limit)
		rate, delay, allowed := limiter.AllowN(factor, limit, time.Second, 0)
		if !allowed {
			c.Header("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
			c.Header("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
			delaySec := int64(delay / time.Second)
			c.Header("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
			c.JSON(ex.ReturnLimitError())
			c.Abort()
		} else {
			c.Next()
		}
	}
}
