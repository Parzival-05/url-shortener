package grpc

import (
	"context"

	url_shortener_v1 "github.com/Parzival-05/url-shortener/api/gen/proto/url_shortener/v1"
	"github.com/Parzival-05/url-shortener/internal/service"
	"go.uber.org/zap"
)

type serverAPI struct {
	url_shortener_v1.UnimplementedUrlShortenerServiceServer

	log          *zap.Logger
	urlShortener service.IUrlShortener
}

func NewServerAPI(log *zap.Logger, urlShortener service.IUrlShortener) url_shortener_v1.UrlShortenerServiceServer {
	return &serverAPI{
		log:          log,
		urlShortener: urlShortener,
	}
}

func (s *serverAPI) CreateShortURL(ctx context.Context, req *url_shortener_v1.CreateShortURLRequest) (*url_shortener_v1.CreateShortURLResponse, error) {
	shortUrl, err := s.urlShortener.CreateUrl(ctx, req.Url)
	return &url_shortener_v1.CreateShortURLResponse{ShortUrl: shortUrl}, err
}

func (s *serverAPI) GetOriginalURL(ctx context.Context, req *url_shortener_v1.GetOriginalURLRequest) (*url_shortener_v1.GetOriginalURLResponse, error) {
	fullUrl, err := s.urlShortener.GetFullUrl(ctx, req.ShortUrl)
	return &url_shortener_v1.GetOriginalURLResponse{Url: fullUrl}, err
}
