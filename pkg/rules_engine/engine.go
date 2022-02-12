package rules_engine

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/effects"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
	json "github.com/json-iterator/go"
	"sync"
)

type RuleDB interface {
	CreateRule(rule Rule) (int64, error)
	UpdateRule(ruleId int64, r RuleDescription) error
	DeleteRule(ruleId int64) error
	GetRules() map[int64]RuleDescription
}

type Engine struct {
	db        RuleDB
	Rules     sync.Map
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
				rule := NewRule(desc.Name, desc.Enabled, tri, eff)
				rule.id = id
				e.Rules.Store(id, rule)
				if r := e.getRule(id); r != nil {
					go r.Start()
				}
			}
		}
		return e.getRules()
	}
	return nil
}

func (e *Engine) GetRule(id int64) []Rule {
	return nil
}

func (e *Engine) AddRule(rule Rule) {

}

func (e *Engine) getRules() map[int64]*Rule {
	rules := make(map[int64]*Rule, 1)
	e.Rules.Range(func(key, value any) bool {
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
	v, ok := e.Rules.Load(id)
	if ok {
		r, ok1 := v.(*Rule)
		if ok1 {
			return r
		}
	}
	return nil
}

func (e *Engine) CreateRule(data []byte) (*Rule, error) {

	var desc RuleDescription
	err := json.Unmarshal(data, &desc)
	if err != nil {
		return nil, err
	}
	rule := FromDescription(desc, e.container)
	id, err := e.db.CreateRule(*rule)
	rule.id = id
	e.Rules.Store(rule.id, rule)

	if err != nil {
		return nil, err
	}
	return rule, nil
}
