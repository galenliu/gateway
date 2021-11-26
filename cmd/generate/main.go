package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

func main() {
	var jsonSchema string
	var cmdStr string
	flag.StringVar(&jsonSchema, "D", "", "DB data path")
	flag.StringVar(&cmdStr, "C", "", "main")
	flag.Parse()
	var file string
	fs, _ := ioutil.ReadDir(path.Join(jsonSchema, "json_schema"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") && !f.IsDir() {
			file = path.Join(jsonSchema, "json_schema", f.Name())
			goFile := path.Join(jsonSchema, "gen", strings.Split(f.Name(), ".json")[0]+".go")
			cmd := exec.Command(
				cmdStr,
				"-p",
				"messages",
				"-o",
				goFile,
				file)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("err : %s \t\n", err.Error())
			}
			fmt.Printf("create file: %s \t\n", goFile)
		} else {
			continue
		}
	}
}
