package log

import (
	"github.com/sirupsen/logrus"
	"io"
)

func Tracef(format string, args ...any) {
	Instance().Tracef(format, args...)
}
func Trace(args ...any) {
	Instance().Trace(args...)
}
func Debugf(format string, args ...any) {
	Instance().Debugf(format, args...)
}
func Debug(args ...any) {
	Instance().Debug(args...)
}
func Infof(format string, args ...any) {
	Instance().Infof(format, args...)
}
func Info(args ...any) {
	Instance().Info(args...)
}
func Warningf(format string, args ...any) {
	Instance().Warningf(format, args...)
}
func Warning(args ...any) {
	Instance().Warning(args...)
}
func Errorf(format string, args ...any) {
	Instance().Errorf(format, args...)
}
func Error(args ...any) {
	Instance().Error(args...)
}
func WithField(key string, value any) *logrus.Entry {
	return Instance().WithField(key, value)
}
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Instance().WithFields(fields)
}
func WriterLevel(lev logrus.Level) *io.PipeWriter {
	return Instance().WriterLevel(lev)
}
func NewEntry() *logrus.Entry {
	return Instance().NewEntry()
}

func Write(p []byte) (n int, err error) {
	Instance().Infof("\u001B[35m" + "web:" + string(p))
	return len(p), nil
}
