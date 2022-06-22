//go:build linux
// +build linux

package cmd

import (
	"errors"
	"github.com/galenliu/gateway/pkg/log"
)

func isWindowsService() (bool, error) {
	return false, nil
}

func createWindowsEventLogger(svcName string, logger log.Logger) (log.Logger, error) {
	return nil, errors.New("cannot create Windows events logger")
}
