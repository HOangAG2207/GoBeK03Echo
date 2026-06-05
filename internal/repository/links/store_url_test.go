package links

import (
	"context"
	"testing"

	pkgredis "github.com/HOangAG2207/GoBeK03Echo/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputCode    string
		inputUrl     string
		inputExpTime int64

		setupMock func() *redis.Client

		expectedError error

		verifyFunc func(ctx context.Context, redisClient *redis.Client)
	}{
		{
			name: "store with no expiration",

			inputCode:    "test",
			inputUrl:     "https://example.com",
			inputExpTime: 0,

			setupMock: func() *redis.Client {
				return pkgredis.InitMockRedis(t)
			},

			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				url, err := redisClient.Get(ctx, "test").Result()

				assert.Nil(t, err)
				assert.Equal(t, "https://example.com", url)
			},
		},
		{
			name: "store with expiration",

			inputCode:    "test",
			inputUrl:     "https://example.com",
			inputExpTime: 60,

			setupMock: func() *redis.Client {
				return pkgredis.InitMockRedis(t)
			},

			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				ttl := redisClient.TTL(ctx, "test").Val()

				assert.Greater(t, ttl.Seconds(), float64(0))
				assert.LessOrEqual(t, ttl.Seconds(), float64(60))
			},
		},
		{
			name: "overwrite existing key",

			inputCode:    "test",
			inputUrl:     "https://example.com",
			inputExpTime: 60,

			setupMock: func() *redis.Client {
				rdb := pkgredis.InitMockRedis(t)
				_ = rdb.Set(context.Background(), "test", "old-value", 0).Err()
				return rdb
			},

			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				url, err := redisClient.Get(ctx, "test").Result()

				assert.Nil(t, err)
				assert.Equal(t, "https://example.com", url)
			},
		},
		{
			name: "failed URL storage due to closed Redis client",

			inputCode:    "test",
			inputUrl:     "https://example.com",
			inputExpTime: 60,

			setupMock: func() *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisMockClient := tc.setupMock()

			urlStorage := NewRepository(redisMockClient)

			err := urlStorage.StoreURL(
				ctx,
				tc.inputCode,
				tc.inputUrl,
				tc.inputExpTime,
			)

			assert.Equal(t, tc.expectedError, err)

			if err == nil && tc.verifyFunc != nil {
				tc.verifyFunc(ctx, redisMockClient)
			}
		})
	}
}
