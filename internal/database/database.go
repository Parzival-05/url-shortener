package database

import "url-shortener/internal/domain"

// DBService represents a service that interacts with a database.
type DBService interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	SyncDB()

	NewUrlRepository() domain.UrlRepository
}
