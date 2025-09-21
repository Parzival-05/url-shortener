package http_server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Parzival-05/url-shortener/internal/database"
	"github.com/Parzival-05/url-shortener/internal/service"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type Server struct {
	port         int
	log          *zap.Logger
	db           database.DBService
	urlShortener service.IUrlShortener
}

func NewServer(log *zap.Logger, db database.DBService, urlShortener service.IUrlShortener) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	NewServer := &Server{
		port:         port,
		log:          log,
		db:           db,
		urlShortener: urlShortener,
	}

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
