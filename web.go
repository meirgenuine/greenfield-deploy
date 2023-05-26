package main

import (
	"context"
	v1 "greenfield-deploy/web/v1"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/mgutz/logxi/v1"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	webCmd = &cobra.Command{
		Use:   "web",
		Short: "Runs web interface for BandLab Geo Internal service",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				router = mux.NewRouter()
				web    = http.Server{
					Addr:    ":8080", // TODO: use viper here
					Handler: router,
				}

				serverContext = context.Background()
				signalChan    = make(chan os.Signal, 1)
			)

			go func() {
				signal.Notify(signalChan, os.Interrupt)
				signal.Notify(signalChan, syscall.SIGTERM)
				log.Info("assigned signal handlers")

				if s, ok := <-signalChan; ok {
					log.Info("signal found", "value", s.String())
				} else {
					log.Info("signal channel closed")
				}
				log.Info("terminating web server...")

				if err := web.Shutdown(serverContext); err != nil {
					log.Error("web server shutdown error", "error", err)
				}
				log.Info("web server terminated")
			}()

			router.HandleFunc("/v1/deploy", v1.DeployHandler)

			log.Info("starting web server", "addr", ":8080")
			if err := web.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal("error on listen and serve", "error", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(webCmd)
}
