package http_server

import (
	"errors"
	"net/http"

	"github.com/Parzival-05/url-shortener/internal/http_server/io_server"
	domain "github.com/Parzival-05/url-shortener/internal/service"

	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

var decoder = schema.NewDecoder()

// @Summary		Create a short URL
// @Description	Creates a new short link for a given URL. If the URL already exists, it returns the existing short link.
// @Tags			URL Shortener
// @Accept			json
// @Produce		json
// @Param			request	body		io_server.CreateUrlRequest	true	"URL to be shortened"
// @Success		200		{object}	io_server.CreateUrlResponse	"Successfully created or retrieved the short URL"
// @Failure		400		{object}	map[string]string			"Bad Request - Invalid JSON format"
// @Failure		500		{object}	map[string]string			"Internal Server Error"
// @Router			/shorten [post]
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
	shortenUrl, err := urlShortener.CreateUrl(ctx, req.URL)
	if err != nil {
		errorResponse(rc,
			ErrorInfo{
				err:      err,
				code:     http.StatusInternalServerError,
				logLevel: zap.ErrorLevel,
			})
		return
	}
	resp := io_server.CreateUrlResponse{
		ShortenURL: shortenUrl,
	}
	okResponse(rc, ResponseInfo{
		code: http.StatusOK,
		data: resp,
	})
}

// @Summary		Get original URL
// @Description	Retrieves the original, full URL for a given short link code.
// @Tags			URL Shortener
// @Produce		json
// @Param			shorten_url	query		string						true	"The 10-character short code"	Format(string)
// @Success		200			{object}	io_server.GetUrlResponse	"Successfully retrieved the original URL"
// @Failure		400			{object}	map[string]string			"Bad Request - The short code is invalid or was not found"
// @Failure		500			{object}	map[string]string			"Internal Server Error"
// @Router			/shorten [get]
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
