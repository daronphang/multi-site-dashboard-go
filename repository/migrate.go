package repository

import (
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
)


func ProvidePgMigrateInstance(driver database.Driver, relativePath string) (*migrate.Migrate, error) {
	cwd, _ := os.Getwd()
	m, err := migrate.NewWithDatabaseInstance(
		path.Join("file:///", cwd, relativePath), 
		"postgres", 
		driver,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}