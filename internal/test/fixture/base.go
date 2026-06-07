package fixture

import (
	"testing"

	"github.com/PhanNam1501/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

type Fixture interface {
	SetupDB(*gorm.DB)
	Migrate() error
	GenerateData() error
	DB() *gorm.DB
}

func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	// create test database
	fix.SetupDB(sqldb.InitMockDB(t))
	// migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatal("Failed to migrate db for testing")
	}
	//create test data
	err = fix.GenerateData()
	if err != nil {
		t.Fatal("Failed to generate data for testing")
	}

	return fix.DB()
}
