package topic

import (
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

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
	ThingPropertyChanged Topic = "thingPropertyChanged"
	ThingEvent           Topic = "thingEvent"
	ThingActionStatus    Topic = "thingActionStatus"

	DeviceAdded           Topic = "deviceAdded"
	DevicePropertyChanged Topic = "devicePropertyChanged"
	DeviceActionStatus    Topic = "deviceActionStatus"
	DeviceEvent           Topic = "deviceEvent"
	DeviceRemoved         Topic = "deviceRemoved"
	DeviceConnected       Topic = "deviceConnected"

	PairingTimeout Topic = "pairingTimeout"
)

type DeviceAddedMessage struct {
	DeviceId string
	devices.Device
}

type DevicePropertyChangedMessage struct {
	DeviceId     string
	PropertyName string
	properties.PropertyDescription
}

type DeviceActionStatusMessage struct {
	DeviceId string
	Action   actions.ActionDescription
}

type ThingAddedMessage struct {
	ThingId string
}

type ThingRemovedMessage struct {
	ThingId string
}

type ThingActionDescription struct {
	Id            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Input         map[string]any `json:"input,omitempty"`
	Status        string         `json:"status,omitempty"`
	TimeRequested string         `json:"timeRequested,omitempty"`
	TimeCompleted string         `json:"timeCompleted,omitempty"`
}
type ThingActionStatusMessage struct {
	ThingId string
	Action  ThingActionDescription
}

type ThingPropertyChangedMessage struct {
	ThingId      string
	PropertyName string
	Value        any
}

type ThingConnectedMessage struct {
	ThingId   string
	Connected bool
}

type ThingModifyMessage struct {
	ThingId string
}

type ThingEventMessage struct {
	ThingId string
}

func (t Topic) ToString() string {
	return string(t)
}

const (
	ValueChanged Topic = "valueChanged"
	StateChanged Topic = "stateChanged"
)
