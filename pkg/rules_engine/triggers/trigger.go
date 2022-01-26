package triggers

type TriggerDescription struct {
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
}

type Trigger struct {
	Type  string
	label string
}

func NewTrigger(des TriggerDescription) *Trigger {
	return &Trigger{
		Type:  des.Type,
		label: des.Label,
	}
}

func (trg *Trigger) ToDescription() *TriggerDescription {
	return &TriggerDescription{
		Type:  trg.Type,
		Label: trg.label,
	}
}

func FromDescription(des TriggerDescription) *Trigger {
	return nil
}
