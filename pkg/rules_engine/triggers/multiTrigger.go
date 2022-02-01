package triggers

import (
	"github.com/galenliu/gateway/pkg/rules_engine"
)

const OpAND = "AND"
const OpOR = "OR"

type MultiTriggerDescription struct {
	Op       string               `json:"op"`
	Triggers []TriggerDescription `json:"triggers"`
}

type MultiTrigger struct {
	*Trigger
	op       string
	triggers []Entity
	id       int64
	states   []bool
	state    bool
}

func NewMultiTrigger(des MultiTriggerDescription) *MultiTrigger {
	m := &MultiTrigger{}
	m.op = des.Op
	m.state = false
	m.triggers = make([]Entity, len(des.Triggers))
	for i, e := range des.Triggers {
		m.triggers[i] = FromDescription(e)
	}
	m.states = make([]bool, len(m.triggers))
	return m
}

func (m *MultiTrigger) ToDescription() *MultiTriggerDescription {
	return nil
}

func (m *MultiTrigger) onStateChanged(triggerIndex int, state State) {
	m.states[triggerIndex] = state.On
	value := m.states[0]
	for i := 1; i < len(m.states); i++ {
		if m.op == OpAND {
			value = value && m.states[i]
		} else if m.op == OpOR {
			value = value || m.states[i]
		}
	}
	if value != m.state {
		m.state = value
		m.Pub(rules_engine.StateChanged, State{
			On: m.state,
		})
	}
}

func (m *MultiTrigger) Start() {
	for i, t := range m.triggers {
		t.Sub(rules_engine.StateChanged, func(state State) {
			m.onStateChanged(i, state)
		})
	}
}

func (m *MultiTrigger) Stop() {

}
