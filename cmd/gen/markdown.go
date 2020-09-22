package main

import (
	"gateway/accessory/gen"
	"gateway/accessory/gen/markdown"
	"os"
	"path"
)

const ResPath = "./accessory/resources"

var AccPath = path.Join(LibPath, "/accessory")
var SevPath = path.Join(AccPath, "/service")
var CharPath = path.Join(AccPath, "/characteristic")

var LibPath, _ = os.Getwd()

func main() {

	gen.ToJson(ResPath, gen.YamlCharFileName, gen.JsonCharFileName)
	markdown.ServiceMarkDown(ResPath, SevPath)
	markdown.AccessoryMarkDown(ResPath, AccPath)
	//resources.ToJson(ResPath, resources.YamlAccFileName, resources.JsonAccFileName)

}
