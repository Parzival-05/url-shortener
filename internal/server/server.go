package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Parzival-05/url-shortener/internal/database"
	"github.com/Parzival-05/url-shortener/internal/database/inmemory"
	"github.com/Parzival-05/url-shortener/internal/database/sql"
	"github.com/Parzival-05/url-shortener/internal/domain"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type IUrlShortener interface {
	GetShortenUrl(ctx context.Context, fullUrl string) (string, error)
	SaveShortenUrl(ctx context.Context, fullUrl string) error
	GetFullUrl(ctx context.Context, shortenUrl string) (string, error)
}

type Server struct {
	port int
	log  *zap.Logger
	db   database.DBService

	urlShortener IUrlShortener
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
	urlShortener := domain.NewUrlShortener(urlRepo, log)

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
