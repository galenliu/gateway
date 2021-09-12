package addon

import (
	"github.com/galenliu/gateway-grpc"
)

type adapter interface {
}

type P interface {
	doPropertyChanged(property *rpc.Property)
}

type Device struct {
	adapter     adapter
	ID          string `json:"id"`
	AtContext   string `json:"@context"`
	AtType      string `json:"@type"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Links               []string `json:"links"`
	BaseHref            string   `json:"baseHref"`
	PinRequired         bool     `json:"pinRequired"`
	CredentialsRequired bool     `json:"credentialsRequired"`

	Properties map[string]P `json:"properties"`
}

func NewDeviceFormString(adapter adapter, id string) *Device {
	device := &Device{}
	device.ID = id
	device.adapter = adapter
	return device
}

func (d *Device) GetID() string {
	return d.ID
}

func (d *Device) GetProperty(name string) P {
	return d.Properties[name]
}

func (d *Device) GetTitle() string {
	return d.Title
}

func (d *Device) GetDescription() string {
	return d.Description
}

func (d *Device) GetAtType() string {
	return d.AtType
}
