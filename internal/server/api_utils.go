package server

import (
	"fmt"
	"net/http"

	"url-shortener/internal/logger/zap_utils"
	"url-shortener/internal/server/io_server"

	"github.com/go-chi/render"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ErrorContext struct {
	w   http.ResponseWriter
	r   *http.Request
	log *zap.Logger
}

type ErrorInfo struct {
	err      error
	code     int
	logLevel zapcore.Level
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

	zapErrMsg := zap_utils.Err(err)
	switch logLevel {
	case zap.InfoLevel:
		log.Info(msg, zapErrMsg)
	case zap.DebugLevel:
		log.Debug(msg, zapErrMsg)
	case zap.WarnLevel:
		log.Warn(msg, zapErrMsg)
	case zap.ErrorLevel:
		log.Error(msg, zapErrMsg)
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
