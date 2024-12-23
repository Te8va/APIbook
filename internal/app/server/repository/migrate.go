package repository

import (
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:generate mockgen -destination=mocks/migrator_mock.gen.go -package=mocks . Migrator
type Migrator interface {
	Up() error
	Close() (sourceErr, databaseErr error)
}

func NewMigrator(filePath string, dsn string) (Migrator, error) {
	return migrate.New(filePath, dsn)
}

func ApplyMigrations(m Migrator) error {

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	sourceErr, databaseErr := m.Close()

	if sourceErr != nil {
		return sourceErr
	}

	if databaseErr != nil {
		return databaseErr
	}

	return nil
}
