package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"runtime"
	"sort"
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
	return runtime.GOOS + "-" + runtime.GOARCH
}

func GetPythonVersion() []string {
	return []string{"3.5", "3,7", "3.8"}
}

func GetNodeVersion() string {
	return "57"
}

func GetDefaultConfigDir() string {
	dir, _ := os.UserHomeDir()
	dirPath := path.Join(dir, constant.ConfDirName)
	return dirPath
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func GetBytes(key interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func InterfaceToFloat64(data interface{}) float64 {
	return ByteToFloat64(GetBytes(data))
}

func JsonIndent(data string) string {
	var bf bytes.Buffer
	err := json.Indent(&bf, []byte(data), "", " ")
	if err != nil {
		return ""
	}
	return bf.String()
}

func IsJson(in []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(in, &js) == nil
}

func In(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}
