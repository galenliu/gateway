package rules_engine

import (
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/effects"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
	"sync"
)

type RuleDB interface {
	CreateRule(rule Rule) (int64, error)
	UpdateRule(ruleId int64, r Rule) error
	DeleteRule(ruleId int64) error
	GetRules() map[int64]RuleDescription
}

type Engine struct {
	db        RuleDB
	rules     sync.Map
	container container.Container
}

func NewEngine(db RuleDB, container container.Container) *Engine {
	e := &Engine{
		db:        db,
		container: container,
	}
	return e
}

func (e *Engine) GetRules() map[int64]*Rule {
	if len(e.getRules()) == 0 {
		ruleMap := e.db.GetRules()
		if len(ruleMap) > 0 {
			for id, desc := range ruleMap {
				tri := triggers.FromDescription(desc.Trigger, e.container)
				eff := effects.FromDescription(desc.Effect, e.container)
				if tri == nil || eff == nil {
					continue
				}
				rule := NewRule(desc.Name, desc.Enabled, tri, eff)
				rule.setId(id)
				e.rules.Store(id, rule)
				if r := e.getRule(id); r != nil {
					go r.Start()
				}
			}
		}
	}
	return e.getRules()
}

func (e *Engine) GetRule(id int64) *Rule {
	return e.getRule(id)
}

func (e *Engine) AddRule(desc RuleDescription) (*Rule, error) {
	rule := FromDescription(desc, e.container)
	id, err := e.db.CreateRule(*rule)
	if err != nil {
		return nil, err
	}
	rule.setId(id)
	e.rules.Store(id, rule)
	rule.Start()
	return rule, nil
}

func (e *Engine) getRules() map[int64]*Rule {
	rules := make(map[int64]*Rule)
	e.rules.Range(func(key, value any) bool {
		id, ok := key.(int64)
		v, ok1 := value.(*Rule)
		if ok && ok1 {
			rules[id] = v
		}
		return true
	})
	return rules
}

func (e *Engine) getRule(id int64) *Rule {
	v, ok := e.rules.Load(id)
	if ok {
		r, ok1 := v.(*Rule)
		if ok1 {
			return r
		}
	}
	return nil
}

func (e *Engine) DeleteRule(i int64) error {
	r := e.getRule(i)
	if r == nil {
		return fmt.Errorf("rule %v does not exist", i)
	}
	r.Stop()
	e.rules.Delete(i)
	err := e.db.DeleteRule(i)
	if err != nil {
		fmt.Printf("delete rule %v failed: %v", i, err)
	}
	return nil
}

func (e *Engine) UpdateRule(i int64, desc RuleDescription) error {
	r := e.getRule(i)
	if r == nil {
		return fmt.Errorf("rule %v does not exist", i)
	}
	rule := FromDescription(desc, e.container)
	if rule == nil {
		return fmt.Errorf("bad messages")
	}
	rule.setId(i)
	err := e.db.UpdateRule(i, *r)
	if err != nil {
		fmt.Printf("delete rule %v failed: %v", i, err)
	}
	r.Stop()
	e.rules.Store(i, rule)
	rule.Start()
	return nil
}
