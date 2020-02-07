package log

import "sync"

/**
* @Author: hexing
* @Date: 19-9-28 上午11:29
 */
const tempName = "core"

var (
	rootLoggerEx Logger
	once         sync.Once
)

type Option struct {
	LogPath    string `json:"logPath"`
	LogName    string `json:"logName"`
	LogLevel   string `json:"logLevel"` //0 debug,1 info,2 warn,3 error,4 fatal
	MaxSize    int    `json:"maxSize"`  //日志分割的尺寸 MB
	MaxBackup  int    `json:"maxBackup"`
	MaxAge     int    `json:"maxAge"` //分割日志保存的时间 day
	JsonFormat bool   `json:"JsonFormat"`
}

func Init(option *Option) {
	once.Do(func() {
		//rootLogLevel := DebugLevel
		logFormat := DefaultLogFormat
		if option == nil {
			option = &Option{
				LogName: tempName,
			}
		}
		rootLoggerEx = newZapRootLoggerEx(option, logFormat)
		//SetLogLevel(rootLoggerEx, rootLogLevel)
	})
}
func GetLoggerEx() Logger {
	return rootLoggerEx
}
func SyncEx() {
	impl, ok := rootLoggerEx.(*zapLoggerImpl)
	if ok {
		_ = impl.mainLogger.Sync()
	}

}
