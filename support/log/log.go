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
	LogPath       string `json:"log_path"`
	LogName       string `json:"log_name"`
	LogLevel      Level  `json:"log_level"` //1 debug,2 info,3 warn,4 error,5 fatal
	MaxSize       int    `json:"max_size"`
	MaxBackup     int    `json:"max_backup"`
	MaxAge        int    `json:"max_age"`
	ConsoleFormat bool   `json:"console_format"`
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
