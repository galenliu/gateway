package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

func EnsureDir(baseDir string, dirs ...string) error {
	_, err := ioutil.ReadDir(baseDir)
	if os.IsNotExist(err) {
		e := os.MkdirAll(baseDir, os.ModePerm)
		if e != nil {
			return e
		}
	}
	for _, dir := range dirs {
		_, err := ioutil.ReadDir(dir)
		if os.IsNotExist(err) {
			e := os.MkdirAll(dir, os.ModePerm)
			if e != nil {
				return e
			}
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

func GetArch() string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	return "linux" + "-" + arch
}

func GetPythonVersion() []string {
	return []string{"3.5", "3,7", "3.8"}
}

func GetNodeVersion() string {
	return "57"
}

func GetDefaultConfigDir() string {
	dir, _ := os.UserHomeDir()
	dirPath := path.Join(dir, ConfDirName)
	return dirPath
}

type Form map[string]string

func NewForm(args ...string) Form {
	m := make(map[string]string, 0)
	for i, _ := range args {
		if i%2 == 0 {
			continue
		}
		m[args[i]] = args[i-1]
	}
	return m
}
