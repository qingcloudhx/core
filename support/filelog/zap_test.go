package flielog

import (
	"go.uber.org/zap"
	"testing"
)

/**
* @Author: hexing
* @Date: 19-9-17 下午12:30
 */

func TestNewZapLogger(t *testing.T) {
	BLogger.Debug("xxxxx", zap.String("key", "value"))
	BLogger.Named("flag").Debug("xxxxxx")
}
