package topic

const (
	SetProperty = "SetProperty"
	GetDevices  = "GetDevices"
	GetThings   = "GetMapOfThings"

	Pair    = "pair"
	Unpair  = "unpair"
	Pending = "pending"
)

type Topic string

const (
	ThingAdded           Topic = "thingAdded"
	ThingRemoved         Topic = "thingRemoved"
	ThingConnected       Topic = "thingConnected"
	ThingModify          Topic = "thingModify"
	ThingEvent           Topic = "thingEvent"
	ThingActionStatus    Topic = "thingActionStatus"
	ThingPropertyChanged Topic = "thingPropertyChanged"

	DeviceAdded           Topic = "deviceAdded"
	DeviceActionStatus    Topic = "deviceActionStatus"
	DeviceEvent           Topic = "deviceEvent"
	DevicePropertyChanged Topic = "devicePropertyChanged"
	DeviceRemoved         Topic = "deviceRemoved"
	DeviceConnected       Topic = "deviceConnected"

	PairingTimeout Topic = "pairingTimeout"
)

func (t Topic) ToString() string {
	return string(t)
}

const (
	ValueChanged Topic = "valueChanged"
	StateChanged Topic = "stateChanged"
)
