package database

import (
	"path"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


func ProvidePgMigrateInstance(driver database.Driver) (*migrate.Migrate, error) {
	_, filename, _, _ := runtime.Caller(0)
	m, err := migrate.NewWithDatabaseInstance(
		path.Join("file:///", path.Dir(filename), "migration"), 
		"postgres", 
		driver,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}