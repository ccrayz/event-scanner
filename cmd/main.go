package main

import (
	"log"

	"github.com/spf13/cobra"

	"ccrayz/event-scanner/cmd/run"
	"ccrayz/event-scanner/config"
	"ccrayz/event-scanner/internal/db"
	"ccrayz/event-scanner/internal/indexer/models"

	"gorm.io/gorm"
)

type App struct {
	IndexerDB *gorm.DB
}

func (a *App) Migrate() {
	if err := models.Migrate(a.IndexerDB); err != nil {
		log.Fatalf("Failed to indexer migrate database: %v", err)
	}
}

var (
	rootCmd = &cobra.Command{
		Use:   "scanner",
		Short: "scanner is a tool to scan events",
	}
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config")
	}

	addDB := db.NewAppDB()
	addDB.InitDB(cfg)
	addDB.Migrate()

	rootCmd.AddCommand(run.NewCommand(addDB))

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
