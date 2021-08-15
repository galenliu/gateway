package internal

type Device struct {
	ID    string `json:"id"`
	Title string `json:"title"`

	Properties map[string]*Property

	AdapterId string `json:"adapterId"`
}

func NewDeviceFormString(des string) *Device {
	return nil
}

func (d *Device) GetId() string {
	return d.ID
}

func (d *Device) GetProperty(name string) *Property {
	prop, ok := d.Properties[name]
	if ok {
		return prop
	}
	return nil
}

func (d *Device) GetPropertyValue(name string) (interface{}, error) {
	return nil, nil
}

func (d *Device) SetConnect(connected bool){

}
