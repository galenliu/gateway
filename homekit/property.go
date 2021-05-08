package homekit

import "github.com/brutella/hc/characteristic"

type Property struct {
	*characteristic.Characteristic
}

func (p *Property) OnPropertyChanged(value interface{}) {

}
