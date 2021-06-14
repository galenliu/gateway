package logging

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
)

var instance *zap.Logger

type Logger interface {
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warningf(format string, args ...interface{})
	Warning(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WriterLevel(logrus.Level) *io.PipeWriter
	NewEntry() *logrus.Entry
}

type logger struct {
	*logrus.Logger
	metrics metrics
}

func New(w io.Writer, level logrus.Level) Logger {
	l := logrus.New()
	l.SetOutput(w)
	l.SetLevel(level)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	metrics := newMetrics()
	l.AddHook(metrics)
	return &logger{
		Logger:  l,
		metrics: metrics,
	}
}








func (l *logger) NewEntry() *logrus.Entry {
	return logrus.NewEntry(l.Logger)
}


func InitLogger(logDir string, debug bool, logRotateDays int) {
	// 设置一些基本日志格式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	hook := lumberjack.Logger{
		Filename:   path.Join(logDir, "log.txt"), // 日志文件路径
		MaxSize:    2,                            // 每个日志文件保存的大小 单位:M
		MaxAge:     logRotateDays,                // 文件最多保存多少天
		MaxBackups: 30,                           // 日志文件最多保存多少个备份
		Compress:   false,                        // 是否压缩
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}

	// 如果是开发环境，同时在控制台上也输出
	if debug {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	// 开启文件及行号
	development := zap.Development()

	// 构造日志
	instance = zap.New(core, caller, development, zap.AddCallerSkip(1))

	Debug("logger init succeed")

}

func Debug(format string, v ...interface{}) {
	instance.Sugar().Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	instance.Sugar().Infof(format, v...)
}

func Error(format string, v ...interface{}) {
	instance.Sugar().Errorf(format, v...)
}
