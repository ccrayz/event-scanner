package main

import (
	"log"

	"github.com/spf13/cobra"

	"ccrayz/event-scanner/cmd/apiserver"
	"ccrayz/event-scanner/config"
	"ccrayz/event-scanner/internal/db"
	"ccrayz/event-scanner/internal/indexer/models"
)

var command *cobra.Command

func main() {
	command = apiserver.NewCommand()

	if err := command.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config")
	}
	db.InitDB(cfg)
	if err := models.Migrate(db.DB); err != nil {
		log.Fatalf("Failed to indexer migrate database: %v", err)
	}
}
