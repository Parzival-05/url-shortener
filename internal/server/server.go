package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Parzival-05/url-shortener/internal/database"
	"github.com/Parzival-05/url-shortener/internal/database/sql"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type Server struct {
	port int
	log  *zap.Logger
	db   database.DBService
}

func NewServer(log *zap.Logger) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		log:  log,
		db:   sql.New(),
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
