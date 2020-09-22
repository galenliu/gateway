package golang

import (
	"bytes"
	"gateway/accessory/gen"
	"text/template"
)

const ServiceTemplate = `//THis File is AUTO-GENERATED
package service{{$Name := .StructName}}
import "smartassistant/accessory/characteristic"

const Type{{.StructName}} = "{{.Type}}"

type {{$Name}} struct {
	*Service{{range .RequiredCharacteristics}}
	{{.}}	*characteristic.{{.}}{{end}}
}

func New{{$Name}}() *{{$Name}} {

	svc := {{$Name}}{}
	svc.Service = New(Type{{.StructName}})
	svc.ServiceName = "{{$Name}}"
{{range .RequiredCharacteristics}}
	svc.{{.}} = characteristic.New{{.}}()
	svc.AddCharacteristics(svc.{{.}}.Characteristic)
{{end}}
	return &svc
}

`

type Service struct {
	Name                    string
	uuid                    string
	RequiredCharacteristics []string
	OptionalCharacteristics []string
}

func (s *Service) StructName() string {
	return s.Name
}

func (s *Service) Type() string {
	return uuidTo(s.uuid)
}

func newService(metadata gen.ServiceMetadata) *Service {
	return &Service{
		Name:                    metadata.Name,
		uuid:                    metadata.UUID,
		RequiredCharacteristics: metadata.RequiredCharacteristics,
		OptionalCharacteristics: metadata.OptionalCharacteristics,
	}
}

func GeneratedServiceCode(in *gen.ServiceMetadata) []byte {
	var out bytes.Buffer
	char := newService(*in)

	t := template.New("sevGo")

	t, err := t.Parse(ServiceTemplate)
	checkError(err)
	err = t.Execute(&out, char)
	checkError(err)
	return out.Bytes()

}
