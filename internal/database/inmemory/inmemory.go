package inmemory

import (
	"context"

	"github.com/Parzival-05/url-shortener/internal/service"
)

type InMemoryDBService struct {
}

func NewInMemoryDBService() *InMemoryDBService {
	return &InMemoryDBService{}
}

func (m *InMemoryDBService) Health() map[string]string {
	stats := map[string]string{}
	stats["db_service"] = "inmemory"
	stats["status"] = "ok"
	return stats
}

func (m *InMemoryDBService) Close() error {
	return nil
}

func (m *InMemoryDBService) SyncDB() {
}

func (m *InMemoryDBService) NewUrlRepository() service.UrlRepository {
	return NewInMemoryUrlRepository()
}

type InMemoryUrlRepository struct {
	urlToId map[string]int64
	idToUrl map[int64]string
}

func NewInMemoryUrlRepository() *InMemoryUrlRepository {
	return &InMemoryUrlRepository{
		urlToId: make(map[string]int64),
		idToUrl: make(map[int64]string),
	}
}

func (m *InMemoryUrlRepository) GetID(ctx context.Context, fullUrl string) (id int64, err error) {
	v, exists := m.urlToId[fullUrl]
	if !exists {
		return 0, service.ErrUrlNotFound
	}
	return v, nil
}

func (m *InMemoryUrlRepository) GetUrlByID(ctx context.Context, id int64) (fullUrl string, err error) {
	v, exists := m.idToUrl[id]
	if !exists {
		return "", service.ErrUrlNotFound
	}
	return v, nil
}

func (m *InMemoryUrlRepository) SaveUrl(ctx context.Context, fullUrl string) (err error) {
	m.urlToId[fullUrl] = int64(len(m.urlToId))
	m.idToUrl[m.urlToId[fullUrl]] = fullUrl
	return nil
}
