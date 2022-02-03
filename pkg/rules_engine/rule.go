package rules_engine

import (
	"github.com/galenliu/gateway/pkg/rules_engine/effects"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
)

type RuleDescription struct {
	Enabled bool
	Trigger triggers.TriggerDescription
	Effect  effects.EffectDescription
	Id      int64
	Name    string
}

type Rule struct {
	id      int64
	name    string
	enabled bool
	trigger triggers.Entity
	effect  effects.Entity
}

func (r *Rule) setId(id int64) {
	r.id = id
}

func (r *Rule) setName(name string) {
	r.name = name
}

func (r *Rule) onTriggerStateChanged(state triggers.State) {
	if !r.enabled {
		return
	}
	r.effect.SetState(state)
}

func (r *Rule) Start() {
	r.trigger.Subscribe(StateChanged, r.onTriggerStateChanged)
	r.trigger.Start()
}

func (r *Rule) Stop() {
	r.trigger.Unsubscribe(StateChanged, r.onTriggerStateChanged)
	r.trigger.Stop()
}

func FromDescription(des RuleDescription) *Rule {
	return nil
}
