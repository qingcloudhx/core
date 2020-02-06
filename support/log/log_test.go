package log

import (
	"testing"
	"time"
)

/**
* @Author: hexing
* @Date: 19-9-28 上午11:33
 */

func TestGetLogerEx(t *testing.T) {
	opt := &Option{
		LogPath:   "",
		LogName:   "hexing",
		LogLevel:  1,
		MaxSize:   10,
		MaxBackup: 10,
	}
	Init(opt)
	for {
		GetLoggerEx().Infof("hxhxh:%s", "xxx")
		GetLoggerEx().Warnf("hxhxh:%s", "xxx")
		time.Sleep(2 * time.Second)
	}
	SyncEx()
}
