//THis File is AUTO-GENERATED
package characteristic

const TypeIdentify = "14"

type Identify struct {
	*Bool
}

func NewIdentify() *Identify {

	char := NewBool(TypeIdentify)
	char.Format = FormatBool
	char.Perms = []string{PermWrite}

	char.SetValue(false)
	return &Identify{char}
}
