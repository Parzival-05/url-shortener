package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Parzival-05/url-shortener/docs"
	"github.com/Parzival-05/url-shortener/internal/server"

	"go.uber.org/zap"
)

var (
	appEnvLocal = "local"
	appEnvProd  = "prod"
)

//	@title			URL Shortener API
//	@version		1.0
//	@description	This is a simple service to shorten URLs.

// @license.name	MIT
// @license.url	https://github.com/Parzival-05/url-shortener/blob/main/LICENSE
func main() {
	storageType := flag.String("storage", string(server.InMemory), fmt.Sprintf("Storage type: '%s' or '%s'", string(server.InMemory), string(server.Postgres)))
	help := flag.Bool("help", false, "help")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	ParseStorageType := func(storageType string) server.StorageType {
		switch storageType {
		case string(server.InMemory):
			return server.InMemory
		case string(server.Postgres):
			return server.Postgres
		default:
			panic("unknown storage type")
		}
	}

	log := setupLogger(os.Getenv("APP_ENV"))
	log.Info("Starting server...", zap.String("app_env", os.Getenv("APP_ENV")), zap.String("storage", *storageType))
	server := server.NewServer(log, ParseStorageType(*storageType))

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Info("Graceful shutdown complete.")
}

func setupLogger(appEnv string) *zap.Logger {
	var logger *zap.Logger
	var err error

	switch appEnv {
	case appEnvLocal:
		logger, err = zap.NewDevelopment()
	case appEnvProd:
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	return logger
}

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
