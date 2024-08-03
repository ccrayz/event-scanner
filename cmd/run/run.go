package run

import (
	"ccrayz/event-scanner/internal/db"
	"ccrayz/event-scanner/internal/indexer"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewCommand(appDB *db.AppDB) *cobra.Command {
	var (
		listenPort       string
		terminateSeconds int
	)
	command := &cobra.Command{
		Use:   "run",
		Short: "Start API the server and indexer",
		Run: func(c *cobra.Command, args []string) {
			apiServer := runAPI(listenPort)
			indexer := runIndexer(appDB)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			gracefulShutdown(apiServer, indexer)
		},
	}
	command.PersistentFlags().StringVar(&listenPort, "port", "8080", "Port to listen on")
	command.PersistentFlags().IntVar(&terminateSeconds, "terminate-seconds", 5, "Seconds to wait before terminating the server")
	return command
}

func runAPI(listenPort string) *http.Server {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiServer := &http.Server{
		Addr:    ":" + listenPort,
		Handler: router.Handler(),
	}

	go func() {
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("API Server started	on port", apiServer.Addr)
	return apiServer
}

func runIndexer(db *db.AppDB) *indexer.Indexer {
	schedule := "@every 2s"
	indexer := indexer.NewIndexer(schedule)
	go func() {
		indexer.Run(db)
		fmt.Println("Indexer started")
	}()

	return indexer
}

func gracefulShutdown(apiServer *http.Server, indexer *indexer.Indexer) {
	stopIndexer(indexer)
	stopAPIServer(apiServer)
}

func stopAPIServer(apiServer *http.Server) {
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("Server exiting")
}

func stopIndexer(indexer *indexer.Indexer) {
	log.Println("Shutdown Indexer ...")
	indexer.Stop()
	log.Println("Indexer exiting")
}
