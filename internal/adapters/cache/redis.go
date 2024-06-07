package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type RedisClient struct {
	client  *redis.Client
	expTime time.Duration
}

func New(redisConn string, expTime time.Duration) (*RedisClient, error) {
	opt, err := redis.ParseURL(redisConn)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return &RedisClient{
		client:  client,
		expTime: expTime,
	}, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) {
	r.client.Set(ctx, key, value, 0)
}

func (r *RedisClient) SetTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	r.client.HSet(ctx, key, value, ttl)
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	bytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)

}
func (r *RedisClient) GetClient() *redis.Client {
	return r.client

}

// GetCachedResponse retrieves the cached response if it exists based on proximity
func (r *RedisClient) GetCachedResponse(query string) ([]byte, bool) {
	ctx := context.Background()
	normalizedQuery := NormalizeText(query)

	// Retrieve all keys (in a real scenario, consider using a better structure for scalability)
	iter := r.client.Scan(ctx, 0, "vector:*", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		storedVector, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var vector string
		json.Unmarshal([]byte(storedVector), &vector)
		similarity := CalculateCosineSimilarity(normalizedQuery, vector)
		if similarity > 0.9 { // Adjust threshold as needed
			cachedResponse, err := r.client.Get(ctx, "response:"+strings.Split(key, ":")[1]).Bytes()
			if err == nil {
				return cachedResponse, true
			}
		}
	}
	return nil, false
}

// SetCachedResponse stores the response in the cache with its vector
func (r *RedisClient) SetCachedResponse(query string, response []byte) {
	ctx := context.Background()
	key := GenerateHash(query)
	normalizedQuery := NormalizeText(query)

	vectorBytes, _ := json.Marshal(normalizedQuery)
	r.client.Set(ctx, "vector:"+key, vectorBytes, r.expTime*time.Millisecond)
	r.client.Set(ctx, "response:"+key, response, r.expTime*time.Millisecond)
}

func SetClinet(db *redis.Client) *RedisClient {

	return &RedisClient{
		client:  db,
		expTime: 1000000,
	}
}
