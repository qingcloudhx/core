package log

import "testing"

/**
* @Author: hexing
* @Date: 19-9-17 下午2:06
 */
func TestNewZapLogger(t *testing.T) {
	log := RootLogger()
	log.Infof("xxxx:%s", "hexing")
}
