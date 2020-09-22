//THis File is AUTO-GENERATED
package characteristic

const TypePairingFeatures = "4F"

type PairingFeatures struct {
	*Int
}

func NewPairingFeatures() *PairingFeatures {

	char := NewInt(TypePairingFeatures)
	char.Format = FormatUint8
	char.Perms = []string{PermRead}

	char.SetValue(0)
	return &PairingFeatures{char}
}
