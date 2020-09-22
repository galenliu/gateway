package gen

import (
	"bytes"
	"encoding/json"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path"
)

const YamlCharFileName = "/yaml/characteristics.yaml"
const YamlSevFileName = "/yaml/services.yaml"
const YamlAccFileName = "/yaml/accessory.yaml"

const JsonCharFileName = "/json/characteristics.json"
const JsonSevFileName = "/json/services.json"
const JsonAccFileName = "/json/accessory.json"

func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		print(err)
	}
	return data
}

func writFile(data []byte, fp string) {
	f, err := os.Create(fp)
	if err != nil {
		print(err)
	}
	defer f.Close()
	_, _ = f.Write(data)

}

func ToYaml(resPath string, jsonFileName string, yamlFileName string) {

	yamlFile := path.Join(resPath, yamlFileName)
	jsonFile := path.Join(resPath, jsonFileName)
	y, _ := yaml.JSONToYAML(readFile(jsonFile))

	writFile(y, yamlFile)
}

func ToJson(resPath string, yamlFileName string, jsonFileName string) {

	var out bytes.Buffer
	yamlFile := path.Join(resPath, yamlFileName)
	jsonFile := path.Join(resPath, jsonFileName)
	d, _ := yaml.YAMLToJSON(readFile(yamlFile))
	_ = json.Indent(&out, d, "", "  ")
	writFile(out.Bytes(), jsonFile)

}

func JsonUnmarshal(in interface{}, resInPath string, filName string) error {
	data := readFile(path.Join(resInPath, filName))
	err := json.Unmarshal(data, in)
	return err
}

func YamlUnmarshal(in interface{}, resInPath string, filName string) error {
	data := readFile(path.Join(resInPath, filName))
	err := yaml.Unmarshal(data, in)
	return err
}
