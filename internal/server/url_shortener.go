package server

import (
	"errors"
	"net/http"

	"github.com/Parzival-05/url-shortener/internal/domain"
	"github.com/Parzival-05/url-shortener/internal/server/io_server"

	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

var decoder = schema.NewDecoder()

func (s *Server) CreateUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rc := RequestContext{
		w:   w,
		r:   r,
		log: s.log,
	}
	urlShortener := s.urlShortener
	var req io_server.CreateUrlRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		errorResponse(rc, ErrorInfo{
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
			errorResponse(rc,
				ErrorInfo{
					err:      err,
					code:     http.StatusInternalServerError,
					logLevel: zap.ErrorLevel,
				})
			return
		}
	} else {
		resp := io_server.CreateUrlResponse{
			ShortenURL: shortenUrl,
		}
		okResponse(rc, ResponseInfo{
			code: http.StatusOK,
			data: resp,
		})
		return
	}
	err = urlShortener.SaveShortenUrl(ctx, req.URL)
	if err != nil {
		errorResponse(rc, ErrorInfo{
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
			errorResponse(rc, ErrorInfo{
				err:      err,
				code:     http.StatusInternalServerError,
				logLevel: zap.ErrorLevel,
			})
			return
		} else {
			errorResponse(rc, ErrorInfo{
				err:      err,
				code:     http.StatusInternalServerError,
				logLevel: zap.ErrorLevel,
				msg:      "Shorten url was not found after saving?!",
			})
			return
		}
	}
	resp := io_server.CreateUrlResponse{
		ShortenURL: shortenUrl,
	}
	okResponse(rc, ResponseInfo{
		code: http.StatusOK,
		data: resp})
}

func (s *Server) GetUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rc := RequestContext{
		w:   w,
		r:   r,
		log: s.log,
	}
	urlShortener := s.urlShortener
	var req io_server.GetUrlRequest
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		errorResponse(rc, ErrorInfo{
			err:      err,
			code:     http.StatusBadRequest,
			logLevel: zap.DebugLevel,
			msg:      "Failed to decode request body",
		})
		return
	}
	fullUrl, err := urlShortener.GetFullUrl(ctx, req.ShortenURL)
	if err != nil {
		if errors.Is(err, domain.ErrUrlNotFound) {
			errorResponse(rc, ErrorInfo{
				err:      err,
				code:     http.StatusBadRequest,
				logLevel: zap.DebugLevel,
			})
		} else {
			errorResponse(rc, ErrorInfo{
				err:      err,
				code:     http.StatusInternalServerError,
				logLevel: zap.ErrorLevel,
				msg:      "Failed to get full url",
			})
		}
		return
	}
	resp := io_server.GetUrlResponse{
		URL: fullUrl,
	}
	okResponse(rc, ResponseInfo{
		code: http.StatusOK,
		data: resp,
	})
}
