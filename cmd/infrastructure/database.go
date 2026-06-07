package infrastructure

import (
	"github.com/PhanNam1501/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

const migrationPath = "file://./migrations"

// InitDatabase initializes database connection and runs migrations
func InitDatabase() (*gorm.DB, error) {
	db, err := sqldb.NewClient("")
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := sqldb.MigartePostgresDB(db, migrationPath, "up", 0); err != nil {
		return nil, err
	}

	return db, nil
}
