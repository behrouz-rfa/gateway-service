package config

import (
	"fmt"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
	"log"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

//type Config struct {
//	ServerAddress string
//	DatabaseURL   string
//	RedisURL      string
//	OpenAIAPIKey  string
//	OpenAIModel   string
//}

type Config struct {
	// Server config
	Server struct {
		Port int    `env:"SERVER_PORT"`
		Host string `env:"SERVER_PORT"`
	}

	// Database config
	Database struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASS"`
		Name     string `env:"DB_NAME"`
	}

	OpenAI struct {
		Key string `env:"OPENAI_API_KEY"`
	}

	NetInfo struct {
		Key string `env:"NET_INFO_KEY"`
	}

	Geo struct {
		Key string `env:"GEO_KEY"`
	}

	Redis struct {
		Host         string `env:"REDIS_HOST"`
		Port         int    `env:"REDIS_PORT"`
		User         string `env:"REDIS_USER"`
		Password     string `env:"REDIS_PASS"`
		Name         string `env:"REDIS_NAME"`
		RedisExpTime int64  `env:"REDIS_EXP_TIME"`
	}
	Jwt struct {
		Secret string `env:"JWT_KEY"`
	}

	Plan struct {
		Credit int `env:"PLAN_CREDIT,default=1000"`
	}
	RateLimit struct {
		Limit string `env:"RATE_LIMIT_PER_HOUR,default=4-H"`
	}
}

func (c *Config) DBConnectionString() string {

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", c.Database.Host, c.Database.User, c.Database.Password, c.Database.Name, c.Database.Port)
}
func (c *Config) RedisConnectionString() string {
	return fmt.Sprintf("redis://%s:%s@%s:%d/%s", c.Redis.User, c.Redis.Password, c.Redis.Host, c.Redis.Port, c.Redis.Name)
}

func LoadAndGet() *Config {
	lg := logger.General.Component("config")
	err := godotenv.Load(".env")

	if err != nil {
		err = godotenv.Load(".env.example")
	}

	if err != nil {
		lg.Info("Error loading .env file")
	}
	var conf = &Config{}

	_, err = env.UnmarshalFromEnviron(conf)

	if err != nil {
		log.Fatal(err)
	}
	return conf
}
