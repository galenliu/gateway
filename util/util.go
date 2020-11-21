package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Info struct {
	Name             string
	SerialNumber     string
	Manufacturer     string
	Model            string
	FirmwareRevision string
	ID               uint64
}

func EnsureDir(baseDir string, dirs ...string) error {
	_, err := ioutil.ReadDir(baseDir)
	if os.IsNotExist(err) {
		ee := os.MkdirAll(baseDir, os.ModePerm)
		if ee != nil {
			return ee
		}
	}
	for _, dir := range dirs {
		_, err := ioutil.ReadDir(dir)
		if os.IsNotExist(err) {
			_ = os.MkdirAll(dir, os.ModePerm)

		}
	}
	return nil
}

func CheckSum(file string, checksum string) bool {

	h := sha256.New()
	f, _ := os.Open(file)
	defer f.Close()
	buf := make([]byte, 1<<20)
	for {
		n, err := io.ReadFull(f, buf)
		if err == nil || err == io.ErrUnexpectedEOF {
			_, err = h.Write(buf[0:n])
			if err != nil {

			}

		} else if err == io.EOF {
			break
		} else {

		}

	}
	r := h.Sum(nil)

	sumCode := fmt.Sprintf("%x", r)

	return sumCode == strings.ToLower(checksum)

}
