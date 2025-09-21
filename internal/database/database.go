package database

import (
	"context"
)

type StorageType string

const (
	InMemory StorageType = "inmemory"
	Postgres StorageType = "postgres"
)

// DBService represents a service that interacts with a database.
type DBService interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	SyncDB()

	NewUrlRepository() IUrlRepository
}

type IUrlRepository interface {
	// GetID returns the ID for a given URL
	GetID(ctx context.Context, fullUrl string) (id int64, err error)
	// GetUrlByID returns the full URL for a given ID
	GetUrlByID(ctx context.Context, id int64) (fullUrl string, err error)
	// SaveUrl saves a new URL
	SaveUrl(ctx context.Context, fullUrl string) (err error)
}
