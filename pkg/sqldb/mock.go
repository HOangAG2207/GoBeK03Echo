package pkgdb

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateTestDB(t *testing.T) *gorm.DB {
	cdn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String())

	db, err := gorm.Open(sqlite.Open(cdn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		t.Fatalf("Fail to create test db: %v", err)
	}

	return db
}
