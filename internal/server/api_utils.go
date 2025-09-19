package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/internal/logger/sl"
	"url-shortener/internal/server/io_server"

	"github.com/go-chi/render"
)

type ErrorContext struct {
	w   http.ResponseWriter
	r   *http.Request
	log *slog.Logger
}

type ErrorInfo struct {
	err      error
	code     int
	logLevel slog.Level
	msg      string
}

func errorResponse(errContext ErrorContext, errInfo ErrorInfo) {
	log := errContext.log
	w := errContext.w
	r := errContext.r

	err := errInfo.err
	code := errInfo.code
	logLevel := errInfo.logLevel
	msg := errInfo.msg

	switch logLevel {
	case slog.LevelInfo:
		log.Info(msg, sl.Err(err))
	case slog.LevelDebug:
		log.Debug(msg, sl.Err(err))
	case slog.LevelWarn:
		log.Warn(msg, sl.Err(err))
	case slog.LevelError:
		log.Error(msg, sl.Err(err))
	}
	w.WriteHeader(code)
	var output string
	if msg == "" {
		output = err.Error()
	} else {
		output = fmt.Sprintf(msg, err.Error())
	}
	render.JSON(w, r, io_server.Error(output))
}
