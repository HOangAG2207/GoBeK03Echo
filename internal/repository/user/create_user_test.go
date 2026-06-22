package user

import (
	"testing"

	"gorm.io/gorm"
)

func TestRepository_CreateUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		setupMockDB func(t *testing.T) *gorm.DB
	}
	testCases := []struct {
		name string
	}{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
		})
	}
}
