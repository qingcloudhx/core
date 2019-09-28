package log

import "testing"

/**
* @Author: hexing
* @Date: 19-9-28 上午11:33
 */

func TestGetLogerEx(t *testing.T) {
	Init("hx")
	GetLogerEx().Infof("hxhxh:%s", "xxx")
	SyncEx()
}
