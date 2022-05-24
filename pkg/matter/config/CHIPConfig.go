package config

import "sync"

const (
	ChipPort    = 5540
	ChipUdcPort = ChipPort + 10
)

type MatterDeviceOptions struct {
	SecuredDevicePort         int
	SecuredCommissionerPort   int
	UnsecuredCommissionerPort int
}

var ins *MatterDeviceOptions
var once sync.Once

func GetInstance() *MatterDeviceOptions {
	once.Do(func() {
		ins = &MatterDeviceOptions{}
	})
	return ins
}
