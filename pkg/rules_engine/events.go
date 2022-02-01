package rules_engine

import "github.com/galenliu/gateway/pkg/bus/topic"

const (
	ValueChanged topic.Topic = "valueChanged"
	StateChanged topic.Topic = "stateChanged"
)
