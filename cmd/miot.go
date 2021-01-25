package main

/*
下载米家所有API接口定义
*/

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gateway/pkg/util"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

//读取Type列表
//读取所有的DeviceType
//http://miot-spec.org/miot-spec-v2/spec/devices
//读取所有的ServiceType
//http://miot-spec.org/miot-spec-v2/spec/services
//读取所有的PropertyType
//http://miot-spec.org/miot-spec-v2/spec/properties
//读取所有的ActionType
//http://miot-spec.org/miot-spec-v2/spec/actions
//读取所有的EventType
//http://miot-spec.org/miot-spec-v2/spec/events

var MiotTypes = map[string]string{
	"devices":   "http://miot-spec.org/miot-spec-v2/spec/devices",
	"service":  "http://miot-spec.org/miot-spec-v2/spec/services",
	"properties": "http://miot-spec.org/miot-spec-v2/spec/properties",
	"actions":   "http://miot-spec.org/miot-spec-v2/spec/actions",
	"events":   "http://miot-spec.org/miot-spec-v2/spec/events",
}

var (
	savaDir string
	PJson   string
)

type types struct {
	Types []string `json:"types"`
}

func init() {
	flag.StringVar(&savaDir, "dir", util.GetDefaultConfigDir(), "download directory")
	flag.StringVar(&PJson, "json", "", "dir")
}

func main() {
	flag.Parse()
	DownloadMiot()
}

//首先下载定义项目
func DownloadMiot() {
	for k, v := range MiotTypes {
		nSlice := strings.Split(v, "/")
		name := nSlice[len(nSlice)-1]
		p := path.Join(savaDir, name)
		util.EnsureDir(p)
		fileName := name + ".json"
		f := path.Join(p, fileName)
		downloadJson(f, v)
		//events not type instance
		if k != "events" {
			downloadInstance(k, f, p)
		}
	}

}

func downloadInstance(typ string, jsonFile string, savePath string) {
	var t types

	file, _ := os.Open(jsonFile)
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	err := json.Unmarshal(data, &t)
	if err != nil {
		fmt.Print(err)
		return
	}
	for _, t := range t.Types {
		sl := strings.SplitAfter(t, ":")
		name := sl[len(sl)-2] + sl[len(sl)-1]
		fn := strings.Replace(sl[len(sl)-2], ":", "", -1) + ".json"
		url := "http://miot-spec.org/miot-spec-v2/spec/" + typ + "?type=urn:miot-spec-v2:" + typ + ":" + name
		f := path.Join(savePath, fn)
		downloadJson(f, url)
	}
}

func downloadJson(file string, url string) {
	fmt.Printf("download url=%v , flie =%v\n", url, file)
	var f *os.File
	f, _ = os.Create(file)
	defer f.Close()
	resp, _ := http.Get(url)
	data, _ := ioutil.ReadAll(resp.Body)
	var d bytes.Buffer
	_ = json.Indent(&d, data, "", "    ")
	_, err := f.Write(d.Bytes())
	if err != nil {
		fmt.Print(err)
	}
}
