package triggers

type TimerTriggerDescription struct {
	*Trigger
	Time      string
	Localized bool
}

type TimerTrigger struct {
	time      string
	localized bool
}

func NewTimerTrigger(des TimerTrigger) *TimerTrigger {
	return nil
}
