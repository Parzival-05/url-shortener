package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Parzival-05/url-shortener/docs"
	"github.com/Parzival-05/url-shortener/internal/database"
	"github.com/Parzival-05/url-shortener/internal/database/inmemory"
	"github.com/Parzival-05/url-shortener/internal/database/sql"
	"github.com/Parzival-05/url-shortener/internal/grpc"
	"github.com/Parzival-05/url-shortener/internal/http_server"
	"github.com/Parzival-05/url-shortener/internal/service"

	"go.uber.org/zap"
	google_grpc "google.golang.org/grpc"
)

var (
	appEnvLocal = "local"
	appEnvProd  = "prod"
)

type ServerType string

const (
	httpServer ServerType = "http"
	grpcServer ServerType = "grpc"
)

type Server interface {
	Shutdown(ctx context.Context) error
}

type HttpServer struct {
	*http.Server
}

type GrpcServer struct {
	*google_grpc.Server
}

func (g *GrpcServer) Shutdown(ctx context.Context) error {
	g.GracefulStop()
	return nil
}

//	@title			URL Shortener API
//	@version		1.0
//	@description	This is a simple service to shorten URLs.

// @license.name	MIT
// @license.url	https://github.com/Parzival-05/url-shortener/blob/main/LICENSE
func main() {
	serverTypeS := flag.String("server", string(httpServer), fmt.Sprintf("Type of server to run: '%s', '%s'", string(httpServer), string(grpcServer)))
	storageTypeS := flag.String("storage", string(database.InMemory), fmt.Sprintf("Storage type: '%s' or '%s'", string(database.InMemory), string(database.Postgres)))
	help := flag.Bool("help", false, "help")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	ParseStorageType := func(storageType string) database.StorageType {
		switch storageType {
		case string(database.InMemory):
			return database.InMemory
		case string(database.Postgres):
			return database.Postgres
		default:
			panic("unknown storage type")
		}
	}

	ParseServerType := func(serverType string) ServerType {
		switch serverType {
		case string(httpServer):
			return httpServer
		case string(grpcServer):
			return grpcServer
		default:
			panic("unknown server type")
		}
	}

	log := setupLogger(os.Getenv("APP_ENV"))
	log.Info("Starting server...", zap.String("app_env", os.Getenv("APP_ENV")), zap.String("storage", *storageTypeS), zap.String("server_type", *serverTypeS))
	storageType := ParseStorageType(*storageTypeS)
	var db database.DBService
	if storageType == database.Postgres {
		db = sql.New()
	} else {
		db = inmemory.NewInMemoryDBService()
	}
	db.SyncDB()
	urlRepo := db.NewUrlRepository()
	urlShortener := service.NewUrlShortener(urlRepo, log)
	serverType := ParseServerType(*serverTypeS)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)
	if serverType == httpServer {
		server := http_server.NewServer(log, db, urlShortener)

		// Run graceful shutdown in a separate goroutine
		go gracefulShutdown(server, done)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	} else {
		grpcApi := grpc.NewServerAPI(log, urlShortener)
		var listener net.Listener
		grpcServer, listener := grpc.New(log, grpcApi)

		// Run graceful shutdown in a separate goroutine
		go gracefulShutdown(&GrpcServer{Server: grpcServer}, done)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("gRPC server failed to start", zap.Error(err))
		}
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

func gracefulShutdown(apiServer Server, done chan bool) {
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
