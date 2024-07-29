package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"ccrayz/event-scanner/cmd/apiserver"
	"ccrayz/event-scanner/config"
	"ccrayz/event-scanner/internal/indexer"
	"ccrayz/event-scanner/internal/indexer/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var command *cobra.Command

type App struct {
	IndexerDB *gorm.DB
}

func (a *App) Migrate() {
	if err := models.Migrate(a.IndexerDB); err != nil {
		log.Fatalf("Failed to indexer migrate database: %v", err)
	}
}

func initDB(cfg *config.Config) *gorm.DB {
	switch cfg.Database.Type {
	case "sqlite":
		DB, err := gorm.Open(sqlite.Open(cfg.Database.Sqlite.Path), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to SQLite database: %v", err)
		}
		return DB
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.User,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.DBName,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.SSLMode,
		)
		DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
		}
		return DB
	default:
		log.Fatalf("Unsupported database type: %s", cfg.Database.Type)
		return nil
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config")
	}

	db := initDB(cfg)
	app := &App{
		IndexerDB: db,
	}
	app.Migrate()

	schedule := "@every 2s"
	log.Println(app.IndexerDB)
	indexer := indexer.NewIndexer(schedule, app.IndexerDB)
	go func() {
		indexer.Run()
		fmt.Println("Indexer started")
	}()
	command = apiserver.NewCommand()

	if err := command.Execute(); err != nil {
		panic(err)
	}

	log.Println("Shutdown Indexer ...")
	indexer.Stop()
	log.Println("Indexer exiting")
}
