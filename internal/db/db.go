package db

import (
	"fmt"
	"log"

	"ccrayz/event-scanner/config"
	"ccrayz/event-scanner/internal/indexer/models"

	"gorm.io/driver/postgres" // Add this line
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AppDB struct {
	DB *gorm.DB
}

func NewAppDB() *AppDB {
	return &AppDB{}
}

// InitDB initializes the database connection
func (a *AppDB) InitDB(cfg *config.Config) {
	var err error

	switch cfg.Database.Type {
	case "sqlite":
		a.DB, err = gorm.Open(sqlite.Open(cfg.Database.Sqlite.Path), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to SQLite database: %v", err)
		}
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.User,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.DBName,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.SSLMode,
		)
		a.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
		}
	default:
		log.Fatalf("Unsupported database type: %s", cfg.Database.Type)
	}
}

func (a *AppDB) Migrate() {
	if err := models.Migrate(a.DB); err != nil {
		log.Fatalf("Failed to indexer migrate database: %v", err)
	}
}
