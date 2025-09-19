package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
	"url-shortener/internal/database"
	"url-shortener/internal/database/sql"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	log  *slog.Logger
	db   database.DBService
}

func NewServer(log *slog.Logger) *http.Server {
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
