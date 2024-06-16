package databases

import (
	"database/sql"
	"log"

	"github.com/LGROW101/assessment-tax/config"
)

var db *sql.DB

func Init(cfg *config.Config) {
	connectionString := cfg.DatabaseURL

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0)
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	return db.Close()
}
