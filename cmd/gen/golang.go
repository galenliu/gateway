package main

import (
	"fmt"
	"gateway/accessory"
	"gateway/accessory/gen"
	"gateway/accessory/gen/golang"
	"os"
	"path"
)

const ResPath = "./accessory/resources"

var LibPath, _ = os.Getwd()

var AccPath = path.Join(LibPath, "/accessory")
var SevPath = path.Join(AccPath, "/service")
var CharPath = path.Join(AccPath, "/characteristic")

func main() {
	GenCharacteristic()
	GenService()

}

func GenCharacteristic() {

	var charMetadata map[string]gen.CharacteristicMetadata

	err := gen.YamlUnmarshal(&charMetadata, ResPath, gen.YamlCharFileName)
	checkError(err)
	var charMetadataUnmarshal []gen.CharacteristicMetadata
	for k, v := range charMetadata {
		v.Name = k
		charMetadataUnmarshal = append(charMetadataUnmarshal, v)
	}
	for _, cm := range charMetadataUnmarshal {
		code := golang.GeneratedCharacteristicCode(&cm)
		saveGoFile(code, CharPath, cm.Name)
	}
}

func GenService() {

	var sevMetadata map[string]gen.ServiceMetadata

	err := gen.YamlUnmarshal(&sevMetadata, ResPath, gen.YamlSevFileName)
	checkError(err)
	var sevMetadataUnmarshal []gen.ServiceMetadata
	for k, v := range sevMetadata {
		v.Name = k
		sevMetadataUnmarshal = append(sevMetadataUnmarshal, v)
	}
	for _, sm := range sevMetadataUnmarshal {
		code := golang.GeneratedServiceCode(&sm)
		saveGoFile(code, SevPath, sm.Name)
	}
}

func saveGoFile(b []byte, p string, fName string) {
	fileName := "/" + accessory.CamelCase(fName) + ".go"
	f, err := os.Create(path.Join(p, fileName))
	checkError(err)
	_, e := f.Write(b)
	checkError(e)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
