package server

import (
	"fmt"
	"net/http"

	"github.com/Parzival-05/url-shortener/internal/logger/zap_utils"
	"github.com/Parzival-05/url-shortener/internal/server/io_server"

	"github.com/go-chi/render"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RequestContext struct {
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

type ResponseInfo struct {
	code int
	data interface{} `json:"data"`
}

func errorResponse(rq RequestContext, errInfo ErrorInfo) {
	log := rq.log
	w := rq.w
	r := rq.r

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

func okResponse(rc RequestContext, ri ResponseInfo) {
	w := rc.w
	r := rc.r

	w.WriteHeader(ri.code)
	render.JSON(w, r, io_server.OK(ri.data))
}
