package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func EnsureDir(log logging.Logger, baseDir string, dirs ...string) {
	_, err := ioutil.ReadDir(baseDir)
	if os.IsNotExist(err) {
		e := os.MkdirAll(baseDir, os.ModePerm)
		if e != nil {
			log.Errorf("create dir %s err: %s", baseDir, err.Error())
		}
		log.Infof("create dir: %s", baseDir)
	} else {
		log.Debugf("dir existed: %s", baseDir)
	}
	for _, dir := range dirs {
		_, err := ioutil.ReadDir(dir)
		if os.IsNotExist(err) {
			e := os.MkdirAll(dir, os.ModePerm)
			if e != nil {
				log.Errorf("create dir %s err: %s", dir, err.Error())
			}
			log.Infof("create dir: %s", dir)
		} else {
			log.Debugf("dir existed: %s", dir)
		}
	}
}

func RemoveDir(dir string) error {
	ff, err := os.Stat(dir)
	if err != nil {
		return err
	} else {
		if ff.IsDir() {
			err = os.RemoveAll(dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckSum(file string, checksum string) bool {
	h := sha256.New()
	f, _ := os.Open(file)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("file colse err: %s", err.Error())
		}
	}(f)
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
	return runtime.GOOS + "-" + runtime.GOARCH
}

func GetPythonVersion() (version []string) {
	cmd := exec.Command("python3", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return
	}
	v := strings.SplitAfter(strings.Split(out.String(), " ")[1], ".")
	version = append(version, strings.TrimSuffix(v[0]+v[1], "."))
	return version
}

func GetNodeVersion() (version string) {
	cmd := exec.Command("node", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return
	}
	v := strings.SplitAfter(out.String(), ".")
	if len(v) < 2 {
		return ""
	}
	return strings.TrimPrefix(strings.TrimSuffix(v[0]+v[1], "."), "v")
}

func JsonIndent(in any) string {
	d, err := json.MarshalIndent(&in, "", "  ")
	if err != nil {
		return ""
	}
	return string(d)
}

func IfThen[T any](b bool, trueValue, falseValue T) T {
	if b {
		return trueValue
	} else {
		return falseValue
	}
}

func Get[T any]() {

}
