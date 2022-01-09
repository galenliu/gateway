package addon

import (
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"time"
)

type PropertyProxy interface {
	GetType() string
	GetTitle() string
	GetName() string
	GetUnit() string
	GetEnum() []any
	GetAtType() string
	GetDescription() string
	GetMinimum() any
	GetMaximum() any
	GetMultipleOf() any
	GetValue() any
	IsReadOnly() bool
	SetValue(a any) bool
	SetTitle(s string) bool
	SetDescription(description string) bool
	ToDescription() PropertyDescription
	ToMessage() messages.Property
}

type DeviceProxy interface {
	GetId() string
	GetAdapter() AdapterProxy
	ToMessage() messages.Device
	GetProperty(id string) PropertyProxy
	SetCredentials(username, password string) error
	NotifyPropertyChanged(p PropertyDescription)
}

type AdapterProxy interface {
	GetId() string
	GetName() string
	GetDevice(deviceId string) DeviceProxy
	SendPropertyChangedNotification(deviceId string, property PropertyDescription)
	Unload()
	CancelPairing()
	StartPairing(timeout time.Duration)
	HandleDeviceSaved(DeviceProxy)
	HandleDeviceRemoved(DeviceProxy)
}

type PropertyLinkElem struct {
}

type PropertyDescription struct {
	Name        *string            `json:"name,omitempty"`
	AtType      *string            `json:"@type,omitempty"`
	Title       *string            `json:"title,omitempty"`
	Type        string             `json:"type,omitempty"`
	Unit        *string            `json:"unit,omitempty"`
	Description *string            `json:"description,omitempty"`
	Minimum     *float64           `json:"minimum,omitempty"`
	Maximum     *float64           `json:"maximum,omitempty"`
	Enum        []any              `json:"enum,omitempty"`
	ReadOnly    *bool              `json:"readOnly,omitempty"`
	MultipleOf  *float64           `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       any                `json:"value,omitempty"`
}
