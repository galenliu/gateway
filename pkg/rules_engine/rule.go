package rules_engine

import (
	"github.com/galenliu/gateway/pkg/rules_engine/effects"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
)

type RuleDescription struct {
}

type Rule struct {
	Id      int64
	Name    string
	Enabled bool
	Trigger triggers.Trigger
	Effect  effects.Effect
}

func (r *Rule) setId(id int64) {
	r.Id = id
}

func (r *Rule) setName(name string) {
	r.Name = name
}

func (r *Rule) onTriggerStateChanged() {

}
