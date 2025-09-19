package zap_utils

import (
	"go.uber.org/zap"
)

func Err(err error) zap.Field {
	return zap.String("error", err.Error())
}
