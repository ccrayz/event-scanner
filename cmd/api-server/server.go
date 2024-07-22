package apiserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var (
		listenPort       string
		terminateSeconds int
	)
	command := &cobra.Command{
		Use:   "api-server",
		Short: "Start API the server",
		Run: func(c *cobra.Command, args []string) {
			router := gin.New()
			router.Use(gin.Logger())
			router.Use(gin.Recovery())

			v1 := router.Group("/api/v1")
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			srv := &http.Server{
				Addr:    ":" + listenPort,
				Handler: router.Handler(),
			}

			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			log.Println("Shutdown Server ...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(terminateSeconds)*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown:", err)
			}

			<-ctx.Done()
			log.Printf("timeout of %d seconds.", terminateSeconds)
			log.Println("Server exiting")
		},
	}
	command.PersistentFlags().StringVar(&listenPort, "port", "8080", "Port to listen on")
	command.PersistentFlags().IntVar(&terminateSeconds, "terminate-seconds", 5, "Seconds to wait before terminating the server")
	return command
}
