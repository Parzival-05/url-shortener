package io_server

type CreateUrlRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type GetUrlRequest struct {
	ShortenURL string `json:"shorten_url" validate:"required"`
}
