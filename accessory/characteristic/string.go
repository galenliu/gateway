package characteristic

type String struct {
	*Characteristic
}

func NewString(typ string) *String {
	char := NewCharacteristic(typ)
	char.Format = FormatString
	return &String{char}
}

// SetValue sets a value
func (c *String) SetValue(str string) {
	c.UpdateValue(str)
}

// GetValue returns the value as string
func (c *String) GetValue() string {
	return c.Characteristic.GetValue().(string)
}

//set the maxlength for string char
func (c *String) SetMax(len int) {
	c.MaxLength = len
}
