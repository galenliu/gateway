package gateway

import (
	"io/ioutil"
	"os"
)

type Info struct {
	Name             string
	SerialNumber     string
	Manufacturer     string
	Model            string
	FirmwareRevision string
	ID               uint64
}

func EnsureConfigPath(baseDir string, dirs ...string) {

	_, err := ioutil.ReadDir(baseDir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(baseDir, os.ModePerm)
	}
	for _, dir := range dirs {
		_, err := ioutil.ReadDir(dir)
		if os.IsNotExist(err) {
			_ = os.MkdirAll(dir, os.ModePerm)
		}
	}
}
