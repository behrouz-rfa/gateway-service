package http

import (
	_ "github.com/behrouz-rfa/gateway-service/docs"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/http/handlers"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/http/middleware"
	"github.com/behrouz-rfa/gateway-service/internal/config"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"log/slog"
	"net/http"
)

// Router is a wrapper for the HTTP router
type Router struct {
	*gin.Engine
}

type Services struct {
	GeoService  ports.GeolocationService
	AiService   ports.OpenAIService
	IpService   ports.IpInfoService
	UserService ports.UserService
	Auth        ports.AuthService
	RedisClient *redis.Client
}

// NewRouter creates a new HTTP router
func NewRouter(
	config *config.Config,
	services Services,
) (*Router, error) {
	router := gin.New()

	// Initialize HTTP handlers
	geoHandler := handlers.NewGeolocationHandler(services.GeoService)
	openAiHandler := handlers.NewOpenAi(services.AiService)
	ipInfoHandler := handlers.NewIpInfoHandler(services.IpService)
	userHandler := handlers.NewUserHandler(services.UserService)

	rateLimiterMiddleware := getRateLimiterMiddleWare(config, services)

	setupPrometheus(router)
	setupMiddleware(router)
	setupSwagger(router)
	defineUserRoutes(router, userHandler, services.Auth)

	defineOpenAIRoutes(router, openAiHandler, services.Auth, rateLimiterMiddleware)
	defineGeolocationRoutes(router, geoHandler, services.Auth, rateLimiterMiddleware)
	defineIPInfoRoutes(router, ipInfoHandler, services.Auth, rateLimiterMiddleware)

	return &Router{router}, nil
}

func getRateLimiterMiddleWare(config *config.Config, services Services) gin.HandlerFunc {
	limiter, err := middleware.NewRateLimiter(services.RedisClient, config.RateLimit.Limit)
	if err != nil {
		return func(context *gin.Context) {
			context.Next()
		}
	}

	// Create a new middleware with the limiter instance.
	return mgin.NewMiddleware(limiter)

}

// setupPrometheus sets up Prometheus metrics
func setupPrometheus(router *gin.Engine) {
	customMetrics := []*ginprometheus.Metric{
		{
			ID:          "test_metric",
			Name:        "test_metric",
			Description: "Counter test metric",
			Type:        "counter",
		},
		{
			ID:          "test_metric_2",
			Name:        "test_metric_2",
			Description: "Summary test metric",
			Type:        "summary",
		},
	}

	p := ginprometheus.NewPrometheus("gin", customMetrics)
	p.Use(router)
}

// setupMiddleware sets up middleware for the router
func setupMiddleware(router *gin.Engine) {
	router.Use(sloggin.New(slog.Default()), gin.Recovery())
}

// setupSwagger sets up Swagger documentation
func setupSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func health(router *gin.Engine) {
	// Define the health check API endpoint
	router.GET("/health", func(c *gin.Context) {

		// If everything is healthy, return a 200 OK response
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Application is healthy",
		})

	})
}

// defineUserRoutes defines the routes for user-related operations
func defineUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler, auth ports.AuthService) {
	v1 := router.Group("/api/v1")
	user := v1.Group("/users")
	{
		user.POST("/register", userHandler.Register)
		user.POST("/login", userHandler.Login)

		authUser := user.Use(middleware.AuthMiddleware(auth))
		{
			authUser.GET("/:id", userHandler.GetUser)

			admin := authUser.Use(middleware.AuthMiddleware(auth))
			{
				admin.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}
}

// defineOpenAIRoutes defines the routes for OpenAI-related operations
func defineOpenAIRoutes(router *gin.Engine, ai *handlers.OpenAi, auth ports.AuthService, rateMiddleware gin.HandlerFunc) {
	v1 := router.Group("/api/v1")
	openai := v1.Group("/openai")
	{
		openai.Use(middleware.AuthMiddleware(auth))
		openai.Use(rateMiddleware)
		openai.POST("", ai.AiRequest)
	}
}

// defineGeolocationRoutes defines the routes for geolocation-related operations
func defineGeolocationRoutes(router *gin.Engine, geoHandler *handlers.GeolocationHandler, auth ports.AuthService, rateMiddleware gin.HandlerFunc) {
	v1 := router.Group("/api/v1")
	geo := v1.Group("/geo")
	{
		geo.Use(rateMiddleware)
		geo.Use(middleware.AuthMiddleware(auth))
		geo.GET("/:ip", geoHandler.GetGeolocation /**/)
	}
}

// defineIPInfoRoutes defines the routes for IP information-related operations
func defineIPInfoRoutes(router *gin.Engine, ipInfoHandler *handlers.IpInfoHandler, auth ports.AuthService, rateMiddleware gin.HandlerFunc) {
	v1 := router.Group("/api/v1")
	ipInfo := v1.Group("/ipinfo")
	{
		ipInfo.Use(rateMiddleware)
		ipInfo.Use(middleware.AuthMiddleware(auth))
		ipInfo.GET("/:ip", ipInfoHandler.GetIpInfo /**/)
	}
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Engine.ServeHTTP(w, req)
}
