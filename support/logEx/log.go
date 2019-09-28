package logEx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

/**
* @Author: hexing
* @Date: 19-9-28 上午11:13
 */

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l *Level) Unmarshal(text string) {
	switch text {
	case "debug", "DEBUG":
		*l = DEBUG
	case "info", "INFO", "": // make the zero value useful
		*l = INFO
	case "warn", "WARN":
		*l = WARN
	case "error", "ERROR":
		*l = ERROR
	case "fatal", "FATAL":
		*l = FATAL
	default:
		*l = INFO
	}
}

type Output int

const (
	FILE Output = iota
	CONSOLE
	ALL
)

type Opt struct {
	LogPath   string
	LogName   string
	LogLevel  Level
	MaxSize   int
	MaxBackup int
	MaxAge    int
	LogOutput Output
}

var (
	BLogger *zap.Logger
)

func Init(name string) {
	BLogger, _ = newZapLogger(
		&Opt{
			LogPath:   os.TempDir(),
			LogName:   name + ".log",
			MaxBackup: 10,
			LogLevel:  DEBUG,
			LogOutput: ALL,
		})
}

func newZapLogger(opt *Opt) (*zap.Logger, io.Writer) {
	var writer io.Writer
	switch opt.LogOutput {
	case FILE:
		writer = newFileWriter(opt)
	case CONSOLE:
		writer = os.Stdout
	case ALL:
		writer = io.MultiWriter(newFileWriter(opt), os.Stdout)
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLevel := asZapLevel(opt.LogLevel)
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(writer), zapLevel)

	if zapLevel == zapcore.DebugLevel {
		return zap.New(core, zap.AddCaller()), writer
	}
	return zap.New(core), writer
}

func newFileWriter(opt *Opt) io.Writer {
	if opt.MaxSize <= 0 {
		opt.MaxSize = 100
	}
	if opt.MaxBackup <= 0 {
		opt.MaxBackup = 10
	}
	if opt.MaxAge <= 0 {
		opt.MaxAge = 28
	}
	return &lumberjack.Logger{
		Filename:   filepath.Join(opt.LogPath, opt.LogName),
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackup,
		MaxAge:     opt.MaxAge,
		LocalTime:  true,
		Compress:   true,
	}
}

func asZapLevel(level Level) zapcore.Level {
	zapLevel := zap.InfoLevel
	switch level {
	case DEBUG:
		zapLevel = zap.DebugLevel
	case INFO:
		zapLevel = zap.InfoLevel
	case WARN:
		zapLevel = zap.WarnLevel
	case ERROR:
		zapLevel = zap.ErrorLevel
	case FATAL:
		zapLevel = zap.FatalLevel
	}
	return zapLevel
}
