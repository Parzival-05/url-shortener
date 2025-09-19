package server

import (
	"errors"
	"net/http"

	"github.com/Parzival-05/url-shortener/internal/domain"
	"github.com/Parzival-05/url-shortener/internal/server/io_server"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func (s *Server) CreateUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	errContext := ErrorContext{
		w:   w,
		r:   r,
		log: s.log,
	}
	resp := make(map[string]string)
	urlRepo := s.db.NewUrlRepository()
	urlShortener := domain.NewUrlShortener(urlRepo, s.log)
	var req io_server.CreateUrlRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		errorResponse(errContext, ErrorInfo{
			err:      err,
			code:     http.StatusBadRequest,
			logLevel: zap.DebugLevel,
			msg:      "Failed to decode request body",
		})
		return
	}
	shortenUrl, err := urlShortener.GetShortenUrl(ctx, req.URL)
	if err != nil {
		if !errors.Is(err, domain.ErrUrlNotFound) {
			errorResponse(errContext,
				ErrorInfo{
					err:      err,
					code:     http.StatusInternalServerError,
					logLevel: zap.ErrorLevel,
				})
			return
		}
	} else {
		resp["shorten_url"] = shortenUrl
		render.JSON(w, r, io_server.OK())
		return
	}
	err = urlShortener.SaveShortenUrl(ctx, req.URL)
	if err != nil {
		errorResponse(errContext, ErrorInfo{
			err:      err,
			code:     http.StatusInternalServerError,
			logLevel: zap.ErrorLevel,
			msg:      "Failed to save shorten url",
		})
		return
	}
	shortenUrl, err = urlShortener.GetShortenUrl(ctx, req.URL)
	if err != nil {
		if !errors.Is(err, domain.ErrUrlNotFound) {
			errorResponse(errContext, ErrorInfo{
				err:      err,
				code:     http.StatusInternalServerError,
				logLevel: zap.ErrorLevel,
			})
			return
		} else {
			errorResponse(errContext, ErrorInfo{
				err:      err,
				code:     http.StatusInternalServerError,
				logLevel: zap.ErrorLevel,
				msg:      "Shorten url was not found after saving?!",
			})
			return
		}
	}
	resp["shorten_url"] = shortenUrl
	render.JSON(w, r, io_server.OK())
}

func (s *Server) GetUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	errContext := ErrorContext{
		w:   w,
		r:   r,
		log: s.log,
	}
	resp := make(map[string]string)
	urlRepo := s.db.NewUrlRepository()
	urlShortener := domain.NewUrlShortener(urlRepo, s.log)
	var req io_server.GetUrlRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		errorResponse(errContext, ErrorInfo{
			err:      err,
			code:     http.StatusBadRequest,
			logLevel: zap.DebugLevel,
			msg:      "Failed to decode request body",
		})
		return
	}
	fullUrl, err := urlShortener.GetFullUrl(ctx, req.ShortenURL)
	if err != nil {
		errorResponse(errContext, ErrorInfo{
			err:      err,
			code:     http.StatusInternalServerError,
			logLevel: zap.ErrorLevel,
			msg:      "Failed to get full url",
		})
		return
	}
	resp["full_url"] = fullUrl
	render.JSON(w, r, io_server.OK())
}
