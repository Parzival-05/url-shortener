package domain

import (
	"context"
	"errors"
	"os"

	"github.com/sqids/sqids-go"
	"go.uber.org/zap"
)

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrInvalidUrl  = errors.New("invalid shorten url")
)

type UrlRepository interface {
	GetID(ctx context.Context, fullUrl string) (id int64, err error)
	GetUrlByID(ctx context.Context, id int64) (fullUrl string, err error)
	SaveUrl(ctx context.Context, fullUrl string) (err error)
}

type UrlShortener struct {
	urlRepo UrlRepository
	log     *zap.Logger
}

func NewUrlShortener(urlRepo UrlRepository, log *zap.Logger) *UrlShortener {
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
