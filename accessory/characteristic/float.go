package characteristic

type Float struct {
	*Characteristic
}

func NewFloat(typ string) *Float {
	char := NewCharacteristic(typ)
	char.Format = FormatFloat
	return &Float{char}
}

//SetValue
func (f *Float) SetValue(value float64) {
	f.UpdateValue(value)

}

func (f *Float) SetMinValue(value float64) {
	f.MinValue = value
}

func (f *Float) SetMaxValue(value float64) {
	f.MaxValue = value

}

func (f *Float) SetStepValue(value float64) {
	f.StepValue = value
}

//GetValue
func (f *Float) GetValue() float64 {
	return f.Characteristic.GetValue().(float64)
}

func (f *Float) GetMinValue() float64 {
	return f.MinValue.(float64)
}

func (f *Float) GetMaxValue() float64 {
	return f.MaxValue.(float64)
}

func (f *Float) GetStepValue() float64 {
	return f.StepValue.(float64)
}
