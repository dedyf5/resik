// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package migrator

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strconv"

	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

func initMigrateInstance() (*migrate.Migrate, error) {
	configREST := config.Load(configEntity.ModuleREST)
	dbConfig := configREST.Database

	if dbConfig.Host == "" || dbConfig.Port == 0 || dbConfig.Schema == "" || dbConfig.Username == "" || dbConfig.Password == "" {
		return nil, fmt.Errorf("error loading complete database configuration")
	}

	migrationFiles, err := fs.Sub(embeddedMigrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("error creating embedded sub-filesystem: %w", err)
	}

	sourceDriver, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return nil, fmt.Errorf("error creating embed source driver: %w", err)
	}

	dbConn, _, err := drivers.NewMySQLConnection(dbConfig, true)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}
	databaseDriver, err := mysql.WithInstance(dbConn, &mysql.Config{})
	if err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("error creating database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("goembed", sourceDriver, dbConfig.Schema, databaseDriver)
	if err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("error creating migrate instance: %w", err)
	}

	return m, nil
}

func RunUp(stepsStr string) error {
	m, err := initMigrateInstance()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	log.Println("Running migrations UP...")
	var migrateErr error

	if stepsStr != "" {
		numSteps, parseErr := strconv.Atoi(stepsStr)
		if parseErr != nil {
			return fmt.Errorf("invalid number of steps: %w", parseErr)
		}
		migrateErr = m.Steps(numSteps)
	} else {
		migrateErr = m.Up()
	}

	if migrateErr != nil && migrateErr != migrate.ErrNoChange {
		return fmt.Errorf("migration UP failed: %w", migrateErr)
	}
	if migrateErr == migrate.ErrNoChange {
		log.Println("Migration UP: No change.")
	} else {
		log.Println("Migration UP completed successfully.")
	}

	return nil
}

func RunDown(stepsStr string) error {
	m, err := initMigrateInstance()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	log.Println("Running migrations DOWN...")
	var migrateErr error

	if stepsStr != "" {
		numSteps, parseErr := strconv.Atoi(stepsStr)
		if parseErr != nil {
			return fmt.Errorf("invalid number of steps: %w", parseErr)
		}
		migrateErr = m.Steps(-numSteps)
	} else {
		migrateErr = m.Down()
	}

	if migrateErr != nil && migrateErr != migrate.ErrNoChange {
		return fmt.Errorf("migration DOWN failed: %w", migrateErr)
	}
	if migrateErr == migrate.ErrNoChange {
		log.Println("Migration DOWN: No change.")
	} else {
		log.Println("Migration DOWN completed successfully.")
	}

	return nil
}

func RunVersion() error {
	m, err := initMigrateInstance()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("error getting migration version: %w", err)
	}
	if err == migrate.ErrNilVersion {
		log.Println("No migration version found.")
	} else {
		log.Printf("Migration version: %d, Dirty: %t", version, dirty)
	}

	return nil
}
