package characteristic

type Bool struct {
	*Characteristic
}

func NewBool(typ string) *Bool {
	bool := NewCharacteristic(typ)
	return &Bool{bool}

}

//SetValue
func (b *Bool) SetValue(value bool) {
	b.UpdateValue(value)

}

//GetValue
func (b *Bool) GetValue() bool {
	return b.Characteristic.GetValue().(bool)
}
