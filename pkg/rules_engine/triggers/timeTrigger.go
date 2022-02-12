package triggers

import (
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TimerTriggerDescription struct {
	TriggerDescription
	Time      string
	Localized bool
}

type TimerTrigger struct {
	*Trigger
	time      string
	localized bool
	timeOut   *time.Timer
	locker    sync.Mutex
}

func NewTimerTrigger(des TimerTriggerDescription) *TimerTrigger {
	return &TimerTrigger{
		Trigger:   NewTrigger(des.TriggerDescription),
		time:      des.Time,
		localized: false,
		locker:    sync.Mutex{},
	}
}

func (t *TimerTrigger) Start() {
	t.scheduleNext()
}

func (t *TimerTrigger) scheduleNext() {
	parts := strings.Split(t.time, ":")
	if len(parts) < 2 {
		return
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, time.Local)
	if now.Before(tm) {
		tm = time.Date(now.Year(), now.Month(), now.Day()+1, hours, minutes, 0, 0, time.Local)
	}
	duration := now.Sub(tm)
	t.Stop()
	t.timeOut = time.AfterFunc(duration, t.SendOn)
}

func (t *TimerTrigger) Stop() {
	t.locker.Lock()
	defer t.locker.Unlock()
	if t.timeOut != nil {
		t.timeOut.Stop()
		t.timeOut = nil
	}
}

func (t *TimerTrigger) SendOn() {
	t.Publish(topic.StateChanged, state.State{
		On:    true,
		Value: time.Now(),
	})
	t.Stop()
	t.timeOut = time.AfterFunc(time.Duration(60*1000), t.SendOff)
}

func (t *TimerTrigger) SendOff() {
	t.Publish(topic.StateChanged, state.State{
		On:    false,
		Value: time.Now(),
	})
	t.scheduleNext()
}
