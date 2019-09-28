package log

import "sync"

/**
* @Author: hexing
* @Date: 19-9-28 上午11:29
 */

var (
	rootLoggerEx Logger
	once         sync.Once
)

func Init(name string) {
	once.Do(func() {
		rootLogLevel := DebugLevel
		logFormat := DefaultLogFormat
		rootLoggerEx = newZapRootLoggerEx(name, logFormat)
		SetLogLevel(rootLoggerEx, rootLogLevel)
	})
}
func GetLogerEx() Logger {
	return rootLoggerEx
}
func SyncEx() {
	impl, ok := rootLoggerEx.(*zapLoggerImpl)

	if ok {
		_ = impl.mainLogger.Sync()
	}

}
