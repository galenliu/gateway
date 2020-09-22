package main

import (
	"gateway/accessory/gen"
)

var jsonInResouse = "/json/manifest.json"
var yamlSavedResouse = "/yaml/manifest.yaml"

func main() {

	gen.ToYaml(ResPath, jsonInResouse, yamlSavedResouse)

}
