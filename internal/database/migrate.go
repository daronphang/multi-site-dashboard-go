package database

import (
	"path"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


func ProvidePgMigrateInstance(driver database.Driver, wd string) (*migrate.Migrate, error) {
	m, err := migrate.NewWithDatabaseInstance(
		path.Join("file:///", wd, "internal/database/migration"), 
		"postgres", 
		driver,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}