package io_server

type CreateUrlRequest struct {
	URL string `json:"url" validate:"required,url" schema:"url"`
}
type CreateUrlResponse struct {
	ShortenURL string `json:"shorten_url" validate:"required" schema:"shorten_url"`
}

type GetUrlRequest struct {
	ShortenURL string `json:"shorten_url" validate:"required" schema:"shorten_url"`
}

type GetUrlResponse struct {
	URL string `json:"url" validate:"required,url" schema:"url"`
}
