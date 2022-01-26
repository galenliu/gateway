package effects

type EffectDescription struct {
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
}

type Effect struct {
	t string
	l string
}

func NewEffect(des EffectDescription) *Effect {
	e := &Effect{
		t: des.Type,
		l: des.Label,
	}
	return e
}

func (e *Effect) ToDescription() *EffectDescription {
	des := &EffectDescription{
		Type:  e.t,
		Label: e.l,
	}
	return des
}
