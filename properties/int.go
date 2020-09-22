package properties

type Int struct {
	*Property

	Minimum int
	Maximum int
	Step    int
	Unit    string
}

func NewInt(t string) *Int {
	var p = NewProperty(t)
	var i = &Int{Property: p}
	return i
}

func (int *Int) SetMinimum(value int) {
	int.Minimum = value
}

func (int *Int) SetMaximum(value int) {
	int.Maximum = value
}

func (int *Int) SetStep(value int) {
	int.Step = value
}

func (int *Int) SetUnit(value int) {
	int.Step = value
}
