//go:build darwin
// +build darwin

package cmd

import (
	"errors"
	"github.com/galenliu/gateway/pkg/logging"
)

func isWindowsService() (bool, error) {
	return false, nil
}

func createWindowsEventLogger(svcName string, logger logging.Logger) (logging.Logger, error) {
	return nil, errors.New("cannot create Windows events logger")
}
