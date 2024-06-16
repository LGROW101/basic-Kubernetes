// migration
package migration

import (
	"log"

	"github.com/LGROW101/assessment-tax/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *config.Config) {
	connectionString := cfg.DatabaseURL

	m, err := migrate.New("file://databases/migrations", connectionString)
	if err != nil {
		log.Fatal("Failed to initialize migration:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migration:", err)
	}

	log.Println("Migration completed")
}
