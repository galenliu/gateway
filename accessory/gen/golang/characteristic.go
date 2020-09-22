package golang

import (
	"bytes"
	"fmt"
	"gateway/accessory/gen"
	"os"
	"strings"
	"text/template"
)

const CharacteristicTemplate = `//THis File is AUTO-GENERATED
package characteristic{{$Name := .StructName}}
{{if .HasValidValues}}
const(
	{{range $k,$v := .ValidValues}}{{$Name}}{{$k | vf}} int = {{$v}}
    {{end}}
){{end}}
const Type{{$Name}} = "{{.Type}}"

type {{$Name}} struct {
	*{{.EmbeddedStructName}}
}

func New{{.StructName}}() *{{.StructName}} {

	char := New{{.EmbeddedStructName}}(Type{{.StructName}})
	char.Format = {{.FormatTypeName}}
	char.Perms = {{.Perms}}
	{{if .HasMinValue}}char.SetMinValue({{.MinValue}}){{end}}
	{{if .HasMaxValue}}char.SetMaxValue({{.MaxValue}}){{end}}
	{{if .HasMinStep}}char.SetStepValue({{.MinStep}}){{end}}
	{{if .HasUnit}}char.Unit = {{.Unit}}{{end}}
	{{if .HasDefaultValue}}char.SetValue({{.DefaultValue}}){{end}}
	return &{{.StructName}}{char}
}
`

var formatType = map[string]string{
	"int":    "FormatInt",
	"uint8":  " FormatUint8",
	"uint16": "FormatUint16",
	"uint32": "FormatUint32",
	"int32":  "FormatInt32",
	"uint64": "FormatUint64",
	"float":  "FormatFloat",
	"bool":   "FormatBool",
	"string": "FormatString",
	"tlv8":   "FormatTLV8",
	"data":   "FormatData",
}

var embeddedStructType = map[string]string{
	"uint8":  "Int",
	"uint16": "Int",
	"uint32": "Int",
	"int32":  "Int",
	"uint64": "Int",
	"int":    "Int",
	"float":  "Float",
	"bool":   "Bool",
	"string": "String",
	"tlv8":   "Bytes",
	"data":   "Data",
}

var unitType = map[string]string{

	"percentage": "UnitPercentage", //百分比
	"arcdegrees": "UnitArcDegrees", //度数
	"celsius":    "UnitCelsius",
	"lux":        "UnitLux",
	"seconds":    "UnitSeconds",
	"ppm":        "UnitPPM",
}

type Characteristic struct {
	structName         string
	Format             string
	FormatTypeName     string
	EmbeddedStructName string

	uuid  string
	perms []string

	Unit string

	MinValue interface{}
	MaxValue interface{}
	MinStep  interface{}

	ValidValues map[string]int
}

func (c *Characteristic) Type() string {
	return uuidTo(c.uuid)
}

func (c *Characteristic) HasMinValue() bool {
	return c.MinValue != nil
}

func (c *Characteristic) HasMaxValue() bool {
	return c.MaxValue != nil
}

func (c *Characteristic) HasMinStep() bool {
	return c.MinStep != nil
}

func (c *Characteristic) HasUnit() bool {
	return c.Unit != ""
}

func (c Characteristic) StructName() string {
	return strings.Replace(c.structName, ".", "_", -1)
}

func (c Characteristic) HasDefaultValue() bool {
	return c.DefaultValue() != nil
}
func (c Characteristic) HasValidValues() bool {
	return c.ValidValues != nil
}

func (c Characteristic) DefaultValue() interface{} {

	switch c.Format {
	case "uint8", "uint16", "uint32", "uint64", "int32", "float", "int":
		if c.HasMinValue() {
			return c.MinValue
		}
		return 0
	case "bool":
		return false
	case "string":
		return `""`
	case "tlv8":
		return "[]byte{}"

	default:
		break
	}
	return nil
}

func (c *Characteristic) Perms() string {

	var perms []string
	for _, perm := range c.perms {
		switch perm {
		case "pr":
			perms = append(perms, "PermRead")
		case "pw":
			perms = append(perms, "PermWrite")

		case "ev":
			perms = append(perms, "PermEvents")

		case "hd":
			perms = append(perms, "PermHidden")
		case "aa":
			perms = append(perms, "PermAdditionalAuthorization")
		default:
			panic(fmt.Sprint("Undefined characteristic perm:%v", perm))

		}
	}
	return "[]string{" + strings.Join(perms, ",") + "}"
}

func newCharacteristic(metadata gen.CharacteristicMetadata) *Characteristic {
	return &Characteristic{

		Format:             metadata.Format,
		FormatTypeName:     formatType[metadata.Format],
		EmbeddedStructName: embeddedStructType[metadata.Format],

		uuid:        metadata.UUID,
		perms:       metadata.Permissions,
		structName:  metadata.Name,
		ValidValues: metadata.ValidValues,

		MinValue: metadata.MinValue,
		MaxValue: metadata.MaxValue,
		MinStep:  metadata.MinStep,
		Unit:     unitType[metadata.Unit],
	}
}

func GeneratedCharacteristicCode(in *gen.CharacteristicMetadata) []byte {
	var out bytes.Buffer
	char := newCharacteristic(*in)

	t := template.New("charGo")
	t = t.Funcs(template.FuncMap{"vf": ValidValueFormat})

	t, err := t.Parse(CharacteristicTemplate)
	checkError(err)
	err = t.Execute(&out, char)
	checkError(err)
	return out.Bytes()

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func uuidTo(i string) string {
	p := strings.Split(i, "-")[0]
	for strings.HasPrefix(p, "0") {
		p = p[1:]
	}
	return p
}

func ValidValueFormat(arg ...interface{}) string {
	var out = ""
	if len(arg) == 1 {
		p := strings.NewReplacer(",", "_", "-", "_", "(", "_", ".", "_", ")", "_")
		out = p.Replace(arg[0].(string))

	}
	return out
}
