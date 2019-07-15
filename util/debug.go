package util

import (
	"github.com/qingcloudhx/core/support/log"
	"runtime"
)

/**
* @Author: hexing
* @Date: 19-7-8 下午6:18
 */
func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	log.RootLogger().Errorf("panic ==> %s\n", string(buf[:n]))
}
