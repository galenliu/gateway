//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"io"
)

func isWindowsService() (bool, error) {
	return svc.IsWindowsService()
}

func createWindowsEventLogger(svcName string, logger log.Logger) (log.Logger, error) {
	el, err := eventlog.Open(svcName)
	if err != nil {
		return nil, err
	}

	winlog := &windowsEventLogger{
		logger: logger,
		winlog: el,
	}

	return winlog, nil
}

type windowsEventLogger struct {
	logger log.Logger
	winlog debug.Log
}

func (l *windowsEventLogger) Tracef(format string, args ...any) {
	// ignore
}

func (l *windowsEventLogger) Trace(args ...any) {
	// ignore
}

func (l *windowsEventLogger) Debugf(format string, args ...any) {
	// ignore
}

func (l *windowsEventLogger) Debug(args ...any) {
	// ignore
}

func (l *windowsEventLogger) Infof(format string, args ...any) {
	_ = l.winlog.Info(1633, fmt.Sprintf(format, args...))
}

func (l *windowsEventLogger) Info(args ...any) {
	_ = l.winlog.Info(1633, fmt.Sprint(args...))
}

func (l *windowsEventLogger) Warningf(format string, args ...any) {
	_ = l.winlog.Warning(1633, fmt.Sprintf(format, args...))
}

func (l *windowsEventLogger) Warning(args ...any) {
	_ = l.winlog.Warning(1633, fmt.Sprint(args...))
}

func (l *windowsEventLogger) Errorf(format string, args ...any) {
	_ = l.winlog.Error(1633, fmt.Sprintf(format, args...))
}

func (l *windowsEventLogger) Error(args ...any) {
	_ = l.winlog.Error(1633, fmt.Sprint(args...))
}

func (l *windowsEventLogger) WithField(key string, value any) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *windowsEventLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *windowsEventLogger) WriterLevel(level logrus.Level) *io.PipeWriter {
	return l.NewEntry().WriterLevel(level)
}

func (l *windowsEventLogger) NewEntry() *logrus.Entry {
	return l.logger.NewEntry()
}

func (l *windowsEventLogger) Write(p []byte) (n int, err error) {
	return l.logger.Write(p)
}
