package middleware

import (
	"github.com/redis/go-redis/v9"
	limiter "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func NewRateLimiter(redis *redis.Client, rateLimit string) (*limiter.Limiter, error) {
	store, err := sredis.NewStoreWithOptions(redis, limiter.StoreOptions{
		Prefix: "limiter_gin_example",
	})
	if err != nil {
		return nil, err
	}

	rate, err := limiter.NewRateFromFormatted(rateLimit)
	if err != nil {
		return nil, err
	}

	return limiter.New(store, rate), nil

}

//// JWTProtected func for specify routes group with JWT authentication.
//// See: https://github.com/gofiber/jwt
//func Authenticate(auth ports.AuthService) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		token, err := util.GetAuthToken(ctx)
//		if err != nil {
//			// Log error - will be reported to Sentry
//			//apiErr.HandleError(ctx, apiErr.ErrForbidden.Detail("auth header is required"))
//			return
//		}
//		//check and get user
//		user, err := auth.GetUser(ctx, token)
//		if err != nil {
//			// Log error - will be reported to Sentry
//			//apiErr.HandleError(ctx, apiErr.ErrUnauthorized.Detail("access denied"))
//			return
//		}
//
//		setUserLastActivity(user.ID)
//
//		ctx.Set(string(common.UserContextKey), user)
//
//		ctx.Next()
//	}
//}
//
//func RateLimit() gin.HandlerFunc {
//	store := memory.NewStore()
//
//	return func(c *gin.Context) {
//		userID := c.GetHeader("X-User-ID")
//		if userID == "" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing User ID"})
//			c.Abort()
//			return
//		}
//
//		var user domain.UserInput
//		if err := db.e.First(&user, userID).Error; err != nil {
//			if err == gorm.ErrRecordNotFound {
//				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//				c.Abort()
//				return
//			}
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//			c.Abort()
//			return
//		}
//
//		var plan models.Plan
//		if err := db.DB.First(&plan, user.PlanID).Error; err != nil {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
//			c.Abort()
//			return
//		}
//
//		rate := fmt.Sprintf("%d-H", plan.RateLimit) // Rate limit per hour
//		rateLimiter := limiter.New(store, limiter.Rate{
//			Period: 1 * time.Hour,
//			Limit:  int64(plan.RateLimit),
//		})
//
//		limiterMiddleware := ginmiddleware.RateLimiter(rateLimiter)
//		limiterMiddleware(c)
//		if c.IsAborted() {
//			return
//		}
//
//		c.Next()
//	}
//}
