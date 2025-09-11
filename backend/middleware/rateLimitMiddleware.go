package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/alimosavifard/zyros-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimitMiddleware(redisClient *redis.Client, rateLimit string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(rateLimit)
	if err != nil {
		utils.InitLogger().Warn().Err(err).Msg("Invalid rate limit format, using default")
		rate = limiter.Rate{Period: time.Hour, Limit: 100}
	}

	store, err := sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix: "zyros_limiter",
	})
	if err != nil {
		utils.InitLogger().Fatal().Err(err).Msg("Failed to create Redis store for limiter")
	}

	limiterInstance := limiter.New(store, rate)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiterCtx, err := limiterInstance.Get(c.Request.Context(), ip)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, "Failed to check rate limit", err)
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(int(limiterCtx.Limit)))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(limiterCtx.Remaining)))
		c.Header("X-RateLimit-Reset", strconv.Itoa(int(limiterCtx.Reset)))

		if limiterCtx.Reached {
			utils.SendError(c, http.StatusTooManyRequests, "Too many requests", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}