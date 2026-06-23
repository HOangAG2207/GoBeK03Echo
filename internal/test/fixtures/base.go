package fixtures_test

import (
	"testing"
	"time"

	pkgdb "github.com/HOangAG2207/GoBeK03Echo/pkg/sqldb"
	"gorm.io/gorm"
)

var TestTime = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

type Fixture interface {
	// setup
	SetupDB(db *gorm.DB)
	// migrate
	Migrate() error
	// generate
	GenerateData() error
	// get db
	DB() *gorm.DB
}

type base struct {
	db *gorm.DB
}

func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

func (b *base) DB() *gorm.DB {
	return b.db
}

func NewFixtureDB(t *testing.T, fix Fixture) *gorm.DB {
	// Setup DB
	fix.SetupDB(pkgdb.CreateTestDB(t))

	// Migrate Schema
	if err := fix.Migrate(); err != nil {
		t.Fatalf("Failed to migrate db for testing")
	}

	// Generate Data
	if err := fix.GenerateData(); err != nil {
		t.Fatalf("Failed to generate test data: %v", err)
	}

	// Return DB
	return fix.DB()
}
