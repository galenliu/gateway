package models

type adapter interface {
}
type Prop interface {
	DoPropertyChanged(property []byte)
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

	Properties map[string]Prop `json:"properties"`
}

func NewDeviceFormString(adapter adapter, id string) *Device {
	device := &Device{}
	device.ID = id
	device.adapter = adapter
	return device
}

func (d *Device) GetId() string {
	return d.ID
}

func (d *Device) GetProperty(name string) Prop {
	return d.Properties[name]
}
