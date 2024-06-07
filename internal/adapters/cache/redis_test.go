package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisClient_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &RedisClient{client: db}

	key := "test-key"
	value := "test-value"

	mock.ExpectSet(key, value, 0).SetVal("OK")

	client.Set(context.Background(), key, value)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisClient_SetTTL(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &RedisClient{client: db}

	key := "test-key"
	value := "test-value"
	ttl := time.Second * 10

	mock.ExpectSet(key, value, ttl).SetVal("OK")

	client.SetTTL(context.Background(), key, value, ttl)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisClient_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &RedisClient{client: db}

	key := "test-key"
	expectedValue := "test-value"
	bytes, _ := json.Marshal(expectedValue)

	mock.ExpectGet(key).SetVal(string(bytes))

	var result string
	err := client.Get(context.Background(), key, &result)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisClient_Get_Error(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &RedisClient{client: db}

	key := "non-existent-key"

	mock.ExpectGet(key).SetErr(redis.Nil)

	var result string
	err := client.Get(context.Background(), key, &result)

	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheMiddleware(t *testing.T) {
	r, err := New("redis://default:foobared@localhost:6379/0", 1000)
	assert.Equal(t, err, nil)
	r.SetCachedResponse("what is go", []byte("go is the best language"))
	v, ok := r.GetCachedResponse("what is go")
	assert.Equal(t, ok, true)
	assert.Equal(t, v, []byte("go is the best language"))

	v, ok = r.GetCachedResponse("Golang in in backend")
	assert.Equal(t, ok, false)
	assert.NotEqual(t, v, []byte("go is the best language"))

	v, ok = r.GetCachedResponse("what is golang")
	assert.Equal(t, ok, true)
	assert.Equal(t, v, []byte("go is the best language"))

}

func TestExpired(t *testing.T) {
	r, err := New("redis://default:foobared@localhost:6379/0", 1)
	assert.Equal(t, err, nil)
	r.SetCachedResponse("what is go", []byte("go is the best language"))
	v, ok := r.GetCachedResponse("what is go")
	assert.Equal(t, ok, false)
	assert.Nil(t, v)

}
