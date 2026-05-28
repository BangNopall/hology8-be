package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	middleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func (m *Middleware) RateLimiter() gin.HandlerFunc {
	var (
		rate = limiter.Rate{
			Period: 1 * time.Hour,
			Limit:  1000,
		}

		store = memory.NewStore()

		instance = limiter.New(store, rate)
	)

	return middleware.NewMiddleware(instance)
}
