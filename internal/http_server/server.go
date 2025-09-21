package http_server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Parzival-05/url-shortener/internal/database"
	"github.com/Parzival-05/url-shortener/internal/database/inmemory"
	"github.com/Parzival-05/url-shortener/internal/database/sql"
	"github.com/Parzival-05/url-shortener/internal/service"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type Server struct {
	port int
	log  *zap.Logger
	db   database.DBService

	urlShortener service.IUrlShortener
}

type StorageType string

const (
	InMemory StorageType = "inmemory"
	Postgres StorageType = "postgres"
)

func NewServer(log *zap.Logger, storageType StorageType) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	var db database.DBService
	if storageType == Postgres {
		db = sql.New()
	} else {
		db = inmemory.NewInMemoryDBService()
	}
	urlRepo := db.NewUrlRepository()
	urlShortener := service.NewUrlShortener(urlRepo, log)

	NewServer := &Server{
		port:         port,
		log:          log,
		db:           db,
		urlShortener: urlShortener,
	}
	NewServer.db.SyncDB()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
