package sqldb

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func MigartePostgresDB(db *gorm.DB, migrationPath string, mode string, steps int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, db.Name(), driver)
	if err != nil {
		return err
	}

	migrateErr := m.Up()
	if migrateErr != nil && !errors.Is(migrateErr, migrate.ErrNoChange) {
		return migrateErr
	}

	return nil
}

func migrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var migrateErr error
	switch mode {
	case "up":
		migrateErr = m.Up()
	case "dowen":
		migrateErr = m.Down()
	case "steps":
		migrateErr = m.Steps(steps)
	default:
		return errors.New("invalid mode")
	}

	if migrateErr != nil || !errors.Is(migrateErr, migrate.ErrNoChange) {
		return migrateErr
	}

	return nil
}
