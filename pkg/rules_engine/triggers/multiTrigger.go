package triggers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	json "github.com/json-iterator/go"
)

const OpAND = "AND"
const OpOR = "OR"

type MultiTriggerDescription struct {
	Op       string `json:"op"`
	Triggers []any  `json:"triggers"`
}

type MultiTrigger struct {
	*Trigger
	op       string
	triggers []Entity
	id       int64
	states   []bool
	state    bool
}

func NewMultiTrigger(des MultiTriggerDescription, container things.Container) *MultiTrigger {
	m := &MultiTrigger{}
	m.op = des.Op
	m.state = false
	m.triggers = make([]Entity, len(des.Triggers))
	for i, e := range des.Triggers {
		m.triggers[i] = FromDescription(e, container)
	}
	m.states = make([]bool, len(m.triggers))
	return m
}

func (m *MultiTrigger) ToDescription() MultiTriggerDescription {
	return MultiTriggerDescription{
		Op: m.op,
		Triggers: func() []any {
			arr := make([]any, 0)
			for _, t := range m.triggers {
				arr = append(arr, t)
			}
			return arr
		}(),
	}
}

func (m *MultiTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToDescription())
}

func (m *MultiTrigger) onStateChanged(triggerIndex int, s state.State) {
	m.states[triggerIndex] = s.On
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
		m.Publish(topic.StateChanged, state.State{
			On: value,
		})
	}
}

func (m *MultiTrigger) Start() {
	for i, t := range m.triggers {
		t.Start()
		t.Subscribe(topic.StateChanged, func(state state.State) {
			m.onStateChanged(i, state)
		})
	}
}

func (m *MultiTrigger) Stop() {
	for i, t := range m.triggers {
		t.Stop()
		t.Unsubscribe(topic.StateChanged, func(state state.State) {
			m.onStateChanged(i, state)
		})
	}
}
