package homekit

type Things struct {
	things map[string]*Thing
}

func NewThings() *Things {
	things := &Things{}
	things.things = make(map[string]*Thing)
	return things
}

func (ts *Things) AddThing(data []byte) {
	thing := NewThing(data)
	ts.things[thing.id] = thing
}
