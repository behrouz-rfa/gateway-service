package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/behrouz-rfa/gateway-service/docs"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/auth"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/cache"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/db"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/geo"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/http"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/ipservice"
	chat "github.com/behrouz-rfa/gateway-service/internal/adapters/openai"
	"github.com/behrouz-rfa/gateway-service/internal/config"
	"github.com/behrouz-rfa/gateway-service/internal/core/services"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
)

// @title                   GatewayService
// @version                 1.0
// @description             This is a simple RESTful service
//
// @contact.name            Behrouz R Faris
// @contact.url             https://github.com/behrouz-rfa/gateway-service
// @contact.email           behrouz-rfa@gmail.com
//
// @license.name            MIT
//
// @BasePath                /api/v1
// @schemes                 http https
//
// @securityDefinitions.apikey BearerAuth
// @in                      header
// @name                    Authorization
// @description             Type "Bearer" followed by a space and the access token.
func main() {
	// Initialize logger
	initLogger()

	// Load configuration
	cfg := config.LoadAndGet()

	// Initialize dependencies
	dbRepo := initDB(cfg)
	authRepo := auth.NewAuth(cfg.Jwt.Secret)
	redisRepo := initCache(cfg)
	gptClient := chat.NewClient(cfg.OpenAI.Key)
	ipRepo := ipservice.NewIPServiceRepo(cfg.Geo.Key)
	geoRepo := geo.NewGeolocationService(cfg.Geo.Key, nil)

	// Initialize services
	userService := services.NewUserService(services.UserServiceOpts{
		UserRepo:   dbRepo,
		Auth:       authRepo,
		PlanCredit: cfg.Plan.Credit,
	})

	aiService := services.NewOpenAIService(services.OpenAIServiceOpts{
		Gpt:       gptClient,
		PlanRepo:  dbRepo,
		UserRepo:  dbRepo,
		RedisRepo: redisRepo,
	})

	ipService := services.NewIpInfoService(services.IpInfoServiceOpts{
		PlanRepo:  dbRepo,
		RedisRepo: redisRepo,
		IpRepo:    ipRepo,
	})

	geoService := services.NewGeolocationService(services.GeolocationServiceOpts{
		PlanRepo:  dbRepo,
		RedisRepo: redisRepo,
		Repo:      geoRepo,
	})

	// Initialize and start HTTP server
	startHTTPServer(cfg, authRepo, redisRepo, userService, aiService, ipService, geoService)
}

func initLogger() {
	args := os.Args[1:]
	disableColors := len(args) > 0 && args[0] == "no-colors"
	logger.Init("main", disableColors)
}

func initDB(cfg *config.Config) *db.DbRepository {
	lg := logger.General.Component("main")
	dbRepo, err := db.NewDB(db.DatabaseConfig{
		Driver: "postgres",
		DSN:    cfg.DBConnectionString(),
	})
	if err != nil {
		lg.WithError(err).Fatal("failed to connect to database")
	}
	return dbRepo
}

func initCache(cfg *config.Config) *cache.RedisClient {
	lg := logger.General.Component("main")
	redisRepo, err := cache.New(cfg.RedisConnectionString(), time.Duration(cfg.Redis.RedisExpTime))
	if err != nil {
		lg.WithError(err).Fatal("failed to connect to redis")
	}
	return redisRepo
}

func startHTTPServer(cfg *config.Config, authRepo *auth.Auth, redisRepo *cache.RedisClient, userService *services.UserService, aiService *services.OpenAIService, ipService *services.IpInfoService, geoService *services.GeolocationService) {
	lg := logger.General.Component("main")
	router, err := http.NewRouter(cfg, http.Services{
		IpService:   ipService,
		GeoService:  geoService,
		UserService: userService,
		AiService:   aiService,
		Auth:        authRepo,
		RedisClient: redisRepo.GetClient(),
	})
	if err != nil {
		lg.WithError(err).Fatal("failed to initialize router")
	}

	listenAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lg.WithField("listen_address", listenAddr).Info("Starting HTTP server")
	if err := router.Serve(listenAddr); err != nil {
		lg.WithError(err).Fatal("failed to start HTTP server")
	}
}
