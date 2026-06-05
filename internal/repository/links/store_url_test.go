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
	type fields struct {
		mockRedis  func() *redis.Client
		verifyFunc func(ctx context.Context, redisClient *redis.Client)
	}
	type args struct {
		inputCode    string
		inputUrl     string
		inputExpTime int64
	}
	testCases := []struct {
		name string

		args args

		fields fields

		expectedError error
	}{
		{
			name: "store with no expiration",
			args: args{
				inputCode:    "test",
				inputUrl:     "https://example.com",
				inputExpTime: 0,
			},
			fields: fields{
				mockRedis: func() *redis.Client {
					return pkgredis.InitMockRedis(t)
				},
				verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
					url, err := redisClient.Get(ctx, "test").Result()

					assert.Nil(t, err)
					assert.Equal(t, "https://example.com", url)
				},
			},

			expectedError: nil,
		},
		{
			name: "store with expiration",
			args: args{
				inputCode:    "test",
				inputUrl:     "https://example.com",
				inputExpTime: 60,
			},
			fields: fields{
				mockRedis: func() *redis.Client {
					return pkgredis.InitMockRedis(t)
				},
				verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
					ttl := redisClient.TTL(ctx, "test").Val()

					assert.Greater(t, ttl.Seconds(), float64(0))
					assert.LessOrEqual(t, ttl.Seconds(), float64(60))
				},
			},

			expectedError: nil,
		},
		{
			name: "overwrite existing key",
			args: args{
				inputCode:    "test",
				inputUrl:     "https://example.com",
				inputExpTime: 60,
			},
			fields: fields{
				mockRedis: func() *redis.Client {
					rdb := pkgredis.InitMockRedis(t)
					_ = rdb.Set(context.Background(), "test", "old-value", 0).Err()
					return rdb
				},
				verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
					url, err := redisClient.Get(ctx, "test").Result()

					assert.Nil(t, err)
					assert.Equal(t, "https://example.com", url)
				},
			},

			expectedError: nil,
		},
		{
			name: "failed URL storage due to closed Redis client",
			args: args{
				inputCode:    "test",
				inputUrl:     "https://example.com",
				inputExpTime: 60,
			},
			fields: fields{
				mockRedis: func() *redis.Client {
					redisClient := pkgredis.InitMockRedis(t)
					redisClient.Close()
					return redisClient
				},
				verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				},
			},

			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisMockClient := tc.fields.mockRedis()

			urlStorage := NewRepository(redisMockClient)

			err := urlStorage.StoreURL(
				ctx,
				tc.args.inputCode,
				tc.args.inputUrl,
				tc.args.inputExpTime,
			)

			assert.Equal(t, tc.expectedError, err)

			if err == nil && tc.fields.verifyFunc != nil {
				tc.fields.verifyFunc(ctx, redisMockClient)
			}
		})
	}
}
