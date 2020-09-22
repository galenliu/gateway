package gen

//全部的Categories
type Metadata struct {
	Categories      *[]CategoryMetadata
	Characteristics *[]CharacteristicMetadata
	Services        *[]ServiceMetadata
}

type CharacteristicMetadata struct {
	UUID   string
	Name   string
	Format string
	Unit   string

	MinValue interface{} `yaml:"minValue,flow"`
	MaxValue interface{} `yaml:"maxValue,flow"`
	MinStep  interface{} `yaml:"minStep,flow"`

	Permissions []string
	ValidValues map[string]int `json:"ValidValues,flow"`
	//ValidBits   []map[interface{}]interface{} `json:"ValidBits,omitempty"`
}

type ServiceMetadata struct {
	UUID string
	Name string

	OptionalCharacteristics []string `yaml:"OptionalCharacteristics"`
	RequiredCharacteristics []string `yaml:"RequiredCharacteristics"`
}

type CategoryMetadata struct {
	Category int
	Name     string
}
