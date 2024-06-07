//#go:build e2e

package gql

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/auth"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/cache"
	db2 "github.com/behrouz-rfa/gateway-service/internal/adapters/db"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/geo"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/http"
	"github.com/behrouz-rfa/gateway-service/internal/adapters/ipservice"
	chat "github.com/behrouz-rfa/gateway-service/internal/adapters/openai"
	"github.com/behrouz-rfa/gateway-service/internal/config"
	"github.com/behrouz-rfa/gateway-service/internal/core/services"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
	"github.com/go-redis/redismock/v9"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestGqlTestSuite(t *testing.T) {
	suite.Run(t, new(GqlTestSuite))
}

type GqlTestSuite struct {
	suite.Suite
	http.Services
	api   *http.Router
	rmock redismock.ClientMock
	db    *gorm.DB
}

const dsn = "file::memory:?cache=shared"

func (s *GqlTestSuite) SetupSuite() {
	logger.Init("main", true)
	lg := logger.General.Component("main")

	//config.Load()

	cfg := &config.Config{}

	dbRepo, err := db2.NewDB(db2.DatabaseConfig{
		Driver: "sqlite",
		DSN:    dsn,
	})
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = dbRepo.GetDb()

	// Initialize authentication service
	authRepo := auth.NewAuth(cfg.Jwt.Secret)

	db, mock := redismock.NewClientMock()
	//redisRepo, err := cache.New(cfg.RedisConnectionString(), time.Duration(cfg.Redis.RedisExpTime))
	redisRepo := cache.SetClinet(db)
	s.rmock = mock

	//if err != nil {
	//	lg.WithError(err).Fatal("failed to connect to redis")
	//	return
	//}
	msg := &openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{Message: openai.ChatCompletionMessage{Content: "Hello, world!"}},
		},
	}

	gptClient := &MockOpenAIClient{
		MockCreateChatCompletion: func(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
			return msg, nil
		},
	}
	client := &chat.Client{AiClient: gptClient}

	//gptClient := chat.NewClient(cfg.OpenAI.Key)
	ipRepo := ipservice.NewIPServiceRepo(cfg.Geo.Key)
	geoRepo := geo.NewGeolocationService(cfg.Geo.Key, nil)

	userService := services.NewUserService(services.UserServiceOpts{
		UserRepo:   dbRepo,
		Auth:       authRepo,
		PlanCredit: 100000000,
	})

	aiService := services.NewOpenAIService(services.OpenAIServiceOpts{
		Gpt:       client,
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

	// Initialize HTTP router
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

	s.api = router
}

func (s *GqlTestSuite) TearDownSuite() {
	dropSQLiteDatabase(s.db)
}
func dropSQLiteDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}

	// Assuming DSN is the file path for SQLite file-based database
	return os.Remove(dsn)
}

// Mock client implementing the OpenAIClient interface
type MockOpenAIClient struct {
	MockCreateChatCompletion func(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error)
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	return m.MockCreateChatCompletion(ctx, req)
}
