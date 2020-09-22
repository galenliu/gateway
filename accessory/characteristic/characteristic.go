package characteristic

import (
	"gateway/logger"
	"net"
)

type GetFunc func() interface{}
type ChangeFunc func(c *Characteristic, newValue, oldValue interface{})
type ConnChangeFunc func(conn net.Conn, c *Characteristic, newValue, oldValue interface{})

type Characteristic struct {
	ID          uint64      `json:"iid"`
	Type        string      `json:"type"`
	Perms       []string    `json:"perms"`
	Description string      `json:"description,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Format      string      `json:"format"`
	Unit        string      `json:"unit,omitempty"`

	Events bool `json:"-"`

	MaxLength int

	MinValue    interface{} `json:"maxValue,omitempty"`
	MaxValue    interface{} `json:"minValue,omitempty"`
	StepValue   interface{} `json:"minStep,omitempty"`
	ValidValues interface{} `json:"validValues,omitempty"`

	valueGetFunc        GetFunc
	valueChangeFunc     []ChangeFunc
	connValueChangeFunc []ConnChangeFunc
}

//public
func NewCharacteristic(typ string) *Characteristic {
	return &Characteristic{
		Type: typ,
	}
}

func (c *Characteristic) GetValue() interface{} {
	return c.getValue(nil)
}

func (c *Characteristic) UpdateValue(value interface{}) {
	c.updateValue(value, nil, false)
}

func (c *Characteristic) GetValueFromConnection(conn net.Conn) interface{} {
	return c.getValue(conn)
}

func (c *Characteristic) UpdateValueFromConnection(v interface{}, conn net.Conn) {
	c.updateValue(v, conn, false)
}

//Private
func (c *Characteristic) getValue(conn net.Conn) interface{} {
	if c.valueGetFunc != nil {
		c.updateValue(c.valueGetFunc(), conn, false)
	}
	return c.Value
}

func (c *Characteristic) isWritable() bool {
	return writePerm(c.Perms)
}

func (c *Characteristic) isReadable() bool {
	return readPerm(c.Perms)
}

func (c *Characteristic) updateValue(value interface{}, conn net.Conn, checkPerms bool) {

	if checkPerms && !c.isWritable() {
		return
	}
	old := c.Value
	if c.isReadable() {
		c.Value = value
		logger.Info.Printf("UpdateValue Old: %v New: %v", old, value)
	}
	if conn != nil {
		return
	} else {
		c.onValueUpdate(c.valueChangeFunc, value, old)
	}

}

func (c *Characteristic) onValueUpdate(funcs []ChangeFunc, newValue, oldValue interface{}) {
	for _, fn := range funcs {
		fn(c, newValue, oldValue)
	}
}

func (c *Characteristic) onValueUpdateFormConn(funcs []ConnChangeFunc, conn net.Conn, newValue, oldValue interface{}) {
	for _, fn := range funcs {
		fn(conn, c, newValue, oldValue)
	}
}

//readPerm returns true whe perm include read permission
func readPerm(perms []string) bool {
	for _, perm := range perms {
		if perm == PermRead {
			return true
		}
	}
	return false
}

//writePerm return true whe perms include write permission
func writePerm(perms []string) bool {
	for _, perm := range perms {
		if perm == PermWrite {
			return true
		}
	}
	return false
}
