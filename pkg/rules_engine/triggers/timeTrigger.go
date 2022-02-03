package triggers

import (
	"github.com/galenliu/gateway/pkg/rules_engine"
	"time"
)

type TimerTriggerDescription struct {
	*TriggerDescription
	Time      string
	Localized bool
}

type TimerTrigger struct {
	*Trigger
	time      string
	localized bool
}

func NewTimerTrigger(des TimerTrigger) *TimerTrigger {
	return nil
}

func (t *TimerTrigger) SendOn() {
	t.Publish(rules_engine.StateChanged, State{
		On:    true,
		Value: time.Now(),
	})
}

func (t *TimerTrigger) SendOff() {

}
