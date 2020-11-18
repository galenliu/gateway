package util

import (
	"io/ioutil"
	"log"
	"path"
)

var configPath = "/Users/liuguilin/Documents/web-things/gateway/config"

func JsonToYaml(jsonFile,yamlFile string){

}

func SavaToConfigDir(fileName string,data []byte){
	file := path.Join(configPath,fileName)
	err :=ioutil.WriteFile(file ,data,777)
	if err != nil{
		log.Print(err)
	}
}