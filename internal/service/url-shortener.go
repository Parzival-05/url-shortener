package service

import (
	"context"
	"errors"
	"os"

	"github.com/Parzival-05/url-shortener/internal/database"
	"github.com/sqids/sqids-go"
	"go.uber.org/zap"
)

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrInvalidUrl  = errors.New("invalid shorten url")
)

type IUrlShortener interface {
	// GetShortenUrl returns the shorten URL for a given full URL
	GetShortenUrl(ctx context.Context, fullUrl string) (string, error)
	// SaveShortenUrl saves a new shorten URL
	SaveShortenUrl(ctx context.Context, fullUrl string) error
	// GetFullUrl returns the full URL for a given shorten URL
	GetFullUrl(ctx context.Context, shortenUrl string) (string, error)
	// CreateUrl creates a new short link for a given URL or returns the existing
	CreateUrl(ctx context.Context, fullUrl string) (string, error)
}

type UrlShortener struct {
	urlRepo database.IUrlRepository
	log     *zap.Logger
}

func NewUrlShortener(urlRepo database.IUrlRepository, log *zap.Logger) *UrlShortener {
	return &UrlShortener{
		urlRepo: urlRepo,
		log:     log,
	}
}

func encodeID(id int64) (string, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SECRET_ALPHABET"),
		MinLength: 10,
	})
	if err != nil {
		return "", err
	}
	shortUrl, err := s.Encode([]uint64{uint64(id)})
	if err != nil {
		return "", err
	}
	return shortUrl, err
}

func decodeToID(str string) (int64, error) {
	var id int64
	s, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SECRET_ALPHABET"),
		MinLength: 10,
	})
	if err != nil {
		return 0, err
	}
	numbers := s.Decode(str)
	if len(numbers) == 0 {
		return 0, ErrInvalidUrl
	}
	id = int64(numbers[0])
	return id, nil
}

func (u *UrlShortener) GetShortenUrl(ctx context.Context, fullUrl string) (string, error) {
	id, err := u.urlRepo.GetID(ctx, fullUrl)
	if err != nil {
		return "", err
	}
	return encodeID(id)
}

func (u *UrlShortener) SaveShortenUrl(ctx context.Context, fullUrl string) error {
	err := u.urlRepo.SaveUrl(ctx, fullUrl)
	return err
}

func (u *UrlShortener) GetFullUrl(ctx context.Context, shortenUrl string) (string, error) {
	id, err := decodeToID(shortenUrl)
	if err != nil {
		return "", err
	}
	return u.urlRepo.GetUrlByID(ctx, id)
}

func (u *UrlShortener) CreateUrl(ctx context.Context, fullUrl string) (string, error) {
	shortenUrl, err := u.GetShortenUrl(ctx, fullUrl)
	if err != nil {
		if !errors.Is(err, ErrUrlNotFound) {
			return "", err
		}
	} else {
		return shortenUrl, nil
	}
	err = u.SaveShortenUrl(ctx, fullUrl)
	if err != nil {
		return "", err
	}
	return u.GetShortenUrl(ctx, fullUrl)
}
