package logging

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"io"
)

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
	Write(p []byte) (n int, err error)
}

type logger struct {
	*logrus.Logger
	metrics metrics
}

func New(w io.Writer, level logrus.Level) Logger {
	l := logrus.New()
	l.SetLevel(level)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}
	l.SetOutput(ansicolor.NewAnsiColorWriter(w))
	//l.SetReportCaller(true)
	metrics := newMetrics()
	l.AddHook(metrics)
	return &logger{
		Logger:  l,
		metrics: metrics,
	}
}

func (l *logger) Write(p []byte) (n int, err error) {
	l.Logger.Infof(string(p))
	return len(p), nil
}

func (l *logger) NewEntry() *logrus.Entry {
	return logrus.NewEntry(l.Logger)
}

func (l *logger) New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
