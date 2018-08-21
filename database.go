package core

import (
	"errors"
	"fmt"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
	"github.com/soupstore/go-core/logging"
)

// MakeBinDataMigration creates a migration source from files packed into the binary with go-bindata
func MakeBinDataMigration(assetNames []string, assetLoader func(name string) ([]byte, error)) *bindata.AssetSource {
	return bindata.Resource(assetNames, assetLoader)
}

// PerformMigration will update the database using migrations in the asset source
func PerformMigration(s *bindata.AssetSource, databaseName string) error {
	sourceDriver, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("postgres://localhost:5432/%s?sslmode=disable", databaseName)
	migration, err := migrate.NewWithSourceInstance("go-bindata", sourceDriver, addr)
	if err != nil {
		return err
	}

	if err := logDatabaseVersion(migration); err != nil {
		return err
	}

	if err = migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logging.Info("Schema up to date")
		} else {
			return err
		}
	} else {
		if err := logDatabaseVersion(migration); err != nil {
			return err
		}
	}

	return nil
}

func logDatabaseVersion(migration *migrate.Migrate) error {
	currentVersion, dirty, err := migration.Version()
	if err != nil {
		switch err.Error() {
		case "no migration":
			logging.Info("No existing schema")
		default:
			return err
		}
	} else {
		logging.Info(fmt.Sprintf("Schema is at version %d", currentVersion))
	}

	if dirty {
		return errors.New("Dirty schema - please manually fix")
	}

	return nil
}
