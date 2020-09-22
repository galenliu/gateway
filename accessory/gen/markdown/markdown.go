package markdown

import (
	"bytes"
	"fmt"
	"gateway/accessory"
	"gateway/accessory/gen"
	"os"
	"text/template"
)

const ServicesTemplate = `| Service | Characteristics | ID |
| ---- | ---- | --- |
{{range $i, $v := .}}| [{{$i }}](./{{$i | cc}}.go) | {{range $v.RequiredCharacteristics}}[{{.}}](../characteristic/{{. | cc}}.go)<br/>{{end}}{{range $v.OptionalCharacteristics}}[{{.}}](../characteristic/{{. | cc}}.go) <small>Optional</small><br/>{{end}} | {{$v.UUID}} |
{{end}}`

const AccessoryTemplate = `| Accessory | Category |
| ---- | ---- |
{{range $i, $v := .}} | [{{$v }}](./{{$v | cc}}.go) | {{$i}} |
{{end}}`

func ServiceMarkDown(resPath string, savePath string) {

	var ServicesMD map[string]gen.ServiceMetadata
	var out bytes.Buffer

	err := gen.JsonUnmarshal(&ServicesMD, resPath, gen.JsonSevFileName)
	checkError(err)

	t := template.New("md")

	t = t.Funcs(template.FuncMap{"cc": camelCase})

	t, err = t.Parse(ServicesTemplate)
	checkError(err)

	err = t.Execute(&out, ServicesMD)
	checkError(err)

	saveMd(out, savePath+"/README.md")

}

func AccessoryMarkDown(resPath string, savePath string) {

	var accMD map[int]string
	var out bytes.Buffer

	err := gen.JsonUnmarshal(&accMD, resPath, gen.JsonAccFileName)
	checkError(err)

	t := template.New("md")

	t = t.Funcs(template.FuncMap{"cc": camelCase})

	t, err = t.Parse(AccessoryTemplate)
	checkError(err)

	err = t.Execute(&out, accMD)
	checkError(err)

	saveMd(out, savePath+"/README.md")

}

func camelCase(args ...interface{}) string {
	var s = ""
	if len(args) == 1 {
		s, _ = args[0].(string)
	}
	return accessory.CamelCase(s)
}
func saveMd(buff bytes.Buffer, fileMD string) {

	f, err := os.Create(fileMD)
	checkError(err)
	_, _ = f.Write(buff.Bytes())
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
