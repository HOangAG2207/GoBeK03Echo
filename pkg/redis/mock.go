package pkgredis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// InitMockRedis khởi tạo một Redis giả lập (in-memory Redis)
// dùng cho unit test, không cần chạy Redis thật
func InitMockRedis(t *testing.T) *redis.Client {

	// miniredis tạo một Redis server giả lập trong memory
	// giúp test nhanh, không phụ thuộc external service
	mock := miniredis.RunT(t)

	// tạo Redis client thật nhưng trỏ vào mock server
	return redis.NewClient(&redis.Options{
		Addr: mock.Addr(), // địa chỉ của Redis fake
	})
}
