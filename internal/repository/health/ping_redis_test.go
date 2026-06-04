package health

import (
	"testing"

	pkgredis "github.com/HOangAG2207/GoBeK03Echo/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_PingRedis(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRedis func() *redis.Client
		expectedError  error
	}{
		{
			name: "ping redis - success",
			setupMockRedis: func() *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				return redisClient
			},
			expectedError: nil,
		},
		{
			name: "ping redis - failure",
			setupMockRedis: func() *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				// đóng kết nối Redis để mô phỏng lỗi
				redisClient.Close()
				return redisClient
			},
			expectedError: redis.ErrClosed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			redisMockClient := tc.setupMockRedis()
			repo := NewRepository(redisMockClient)

			err := repo.PingRedis(ctx)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
