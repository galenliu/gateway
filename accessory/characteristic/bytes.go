package characteristic

type Bytes struct {
	*Characteristic
}

func NewBytes(typ string) *Bytes {

	char := NewCharacteristic(typ)
	return &Bytes{char}
}

//SetValue
func (f *Bytes) SetValue(value []byte) {
	f.UpdateValue(value)

}

func (f *Bytes) SetMinValue(value []byte) {
	f.MinValue = value
}

func (f *Bytes) SetMaxValue(value []byte) {
	f.MaxValue = value

}

func (f *Bytes) SetStepValue(value []byte) {
	f.StepValue = value
}

//GetValue
func (f *Bytes) GetValue() []byte {
	return f.Characteristic.GetValue().([]byte)
}

func (f *Bytes) GetMinValue() []byte {
	return f.MinValue.([]byte)
}

func (f *Bytes) GetMaxValue() []byte {
	return f.MaxValue.([]byte)
}

func (f *Bytes) GetStepValue() []byte {
	return f.StepValue.([]byte)
}
