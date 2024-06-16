package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	cfg "github.com/nishanthrk/aspire-lms/app/configs"
	"github.com/nishanthrk/aspire-lms/app/logger"
)

func RunMigrations() error {
	db, _ := sql.Open("mysql", cfg.GetConfig().Mysql.GetMysqlConnectionForMigrate())
	driver, _ := mysql.WithInstance(db, &mysql.Config{})

	if err := applyMigrations(driver, "file://migrations/schema"); err != nil {
		return fmt.Errorf("could not apply product migrations: (%v)", err)
	}

	return nil
}

func applyMigrations(driver database.Driver, migrationsPath string) error {
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %v", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		logger.Sugar.Info("No new migrations to apply")
	}

	return nil
}
