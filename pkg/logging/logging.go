package logging

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"io"
)

type Logger interface {
	Tracef(format string, args ...any)
	Trace(args ...any)
	Debugf(format string, args ...any)
	Debug(args ...any)
	Infof(format string, args ...any)
	Info(args ...any)
	Warningf(format string, args ...any)
	Warning(args ...any)
	Errorf(format string, args ...any)
	Error(args ...any)
	WithField(key string, value any) *logrus.Entry
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
	l.Logger.Infof("\u001B[35m" + "web:" + string(p))
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
