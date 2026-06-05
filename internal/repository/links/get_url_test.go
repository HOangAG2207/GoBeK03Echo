package links

import (
	"context"
	"testing"

	pkgredis "github.com/HOangAG2207/GoBeK03Echo/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetURL(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockRedis  func() *redis.Client
		verifyFunc func(ctx context.Context, redis *redis.Client)
	}
	testCases := []struct {
		name string

		fields fields

		expectedError error
	}{

		{
			name: "Get Url success",
			fields: fields{
				mockRedis: func() *redis.Client {
					redisClient := pkgredis.InitMockRedis(t)
					redisClient.Set(t.Context(), "test", "https://example.com", 10000)
					return redisClient
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
			name: "failed Get URL storage due to closed Redis client",
			fields: fields{
				mockRedis: func() *redis.Client {
					redisClient := pkgredis.InitMockRedis(t)
					redisClient.Close()
					return redisClient
				},
				verifyFunc: func(ctx context.Context, redis *redis.Client) {},
			},

			expectedError: redis.ErrClosed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockRedis := tc.fields.mockRedis()

			repo := NewRepository(mockRedis)

			_, err := repo.GetURL(ctx, "test")
			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				tc.fields.verifyFunc(ctx, mockRedis)
			}
		})
	}
}
