package rules_engine

import "sync"

type RuleDB interface {
	CreateRule(description RuleDescription) int64
	UpdateRule(ruleId int64, r Rule) error
	DeleteRule(ruleId int64)
	GetRules() map[int64]RuleDescription
}

type Engine struct {
	db    RuleDB
	Rules sync.Map
}

func (e *Engine) GetRules() map[int64]*Rule {
	if len(e.getRules()) == 0 {
		ruleDesMap := e.db.GetRules()
		if len(ruleDesMap) > 0 {
			for id, ruleDes := range ruleDesMap {
				ruleDes.Id = id
				rule := FromDescription(ruleDes)
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
