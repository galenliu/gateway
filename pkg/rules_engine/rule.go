package rules_engine

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/rules_engine/effects"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
	json "github.com/json-iterator/go"
)

type RuleDescription struct {
	Enabled bool   `json:"enabled"`
	Trigger any    `json:"trigger"`
	Effect  any    `json:"effect"`
	Id      int64  `json:"id"`
	Name    string `json:"name"`
}

type Rule struct {
	id      int64
	name    string
	enabled bool
	trigger triggers.Entity
	effect  effects.Entity
}

func FromDescription(desc RuleDescription, container things.Container) *Rule {
	tri := triggers.FromDescription(desc.Trigger, container)
	eff := effects.FromDescription(desc.Effect, container)
	enabled := desc.Enabled
	name := desc.Name
	return NewRule(name, enabled, tri, eff)
}

func NewRule(name string, enabled bool, trigger triggers.Entity, effect effects.Entity) *Rule {
	return &Rule{
		name:    name,
		enabled: enabled,
		trigger: trigger,
		effect:  effect,
	}
}

func (r *Rule) setId(id int64) {
	r.id = id
}

func (r *Rule) setName(name string) {
	r.name = name
}

func (r Rule) MarshalJSON() ([]byte, error) {
	desc := RuleDescription{
		Enabled: r.enabled,
		Id:      r.id,
		Trigger: func() json.Any {
			data, _ := json.Marshal(&r.trigger)
			return json.Get(data)
		}(),
		Effect: func() json.Any {
			data, _ := json.Marshal(&r.effect)
			return json.Get(data)
		}(),
		Name: r.name,
	}
	return json.Marshal(desc)
}

func (r *Rule) onTriggerStateChanged(state state.State) {
	if !r.enabled {
		return
	}
	r.effect.SetState(state)
}

func (r *Rule) Start() {
	r.trigger.Subscribe(topic.StateChanged, r.onTriggerStateChanged)
	r.trigger.Start()
}

func (r *Rule) Stop() {
	r.trigger.Unsubscribe(topic.StateChanged, r.onTriggerStateChanged)
	r.trigger.Stop()
}
