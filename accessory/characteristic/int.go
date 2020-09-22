package characteristic

type Int struct {
	*Characteristic
}

func NewInt(typ string) *Int {
	char := NewCharacteristic(typ)
	return &Int{char}
}

//SetValue
func (i *Int) SetValue(value int) {
	i.UpdateValue(value)
}

func (i *Int) SetMinValue(value int) {
	i.MinValue = value
}

func (i *Int) SetMaxValue(value int) {
	i.MaxValue = value

}

func (i *Int) SetStepValue(value int) {
	i.StepValue = value
}

//GetValue
func (i *Int) GetValue() int {
	return i.Characteristic.GetValue().(int)
}

func (i *Int) GetMinValue() int {
	return i.MinValue.(int)
}

func (i *Int) GetMaxValue() int {
	return i.MaxValue.(int)
}

func (i *Int) GetStepValue() int {
	return i.StepValue.(int)
}
