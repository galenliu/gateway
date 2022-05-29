package platform

import (
	"github.com/galenliu/gateway/pkg/matter/device"
	"net"
	"strings"
	"sync"
)

type Config struct {
	ChipDeviceConfigPairingSecondaryHint        int
	ChipDeviceConfigDeviceVendorName            string
	ChipDeviceConfigDeviceType                  int
	ChipDeviceConfigDeviceProductName           string
	ChipDeviceConfigPairingInitialHint          int
	ChipDeviceConfigPairingInitialInstruction   string
	ChipDeviceConfigPairingSecondaryInstruction string
}

var cmInstance *ConfigurationManager
var cmOnce sync.Once

func ConfigurationMgr() *ConfigurationManager {
	cmOnce.Do(func() {
		cmInstance = newConfigurationManager(Config{
			ChipDeviceConfigPairingSecondaryHint:        0,
			ChipDeviceConfigDeviceVendorName:            "",
			ChipDeviceConfigDeviceType:                  0,
			ChipDeviceConfigDeviceProductName:           "",
			ChipDeviceConfigPairingInitialHint:          0,
			ChipDeviceConfigPairingInitialInstruction:   "",
			ChipDeviceConfigPairingSecondaryInstruction: "",
		})
	})
	return cmInstance
}

type ConfigurationManager struct {
	mVendorId                                  int
	mVendorName                                string
	mProductName                               string
	mProductId                                 int
	mDeviceType                                device.MatterDeviceType
	mDeviceName                                string
	mTcpSupported                              bool
	mDevicePairingHint                         int
	mDevicePairingSecondaryHint                int
	mDeviceSecondaryPairingHint                int
	mDeviceConfigPairingInitialInstruction     string
	mDeviceConfigPairingSecondaryInstruction   string
	deviceConfigEnableCommissionableDeviceType bool
}

func newConfigurationManager(conf Config) *ConfigurationManager {
	return &ConfigurationManager{
		mVendorId:                                0,
		mVendorName:                              "",
		mProductName:                             "",
		mProductId:                               0,
		mDeviceType:                              conf.ChipDeviceConfigDeviceType,
		mDeviceName:                              conf.ChipDeviceConfigDeviceProductName,
		mDevicePairingHint:                       conf.ChipDeviceConfigPairingInitialHint,
		mDeviceConfigPairingInitialInstruction:   conf.ChipDeviceConfigPairingInitialInstruction,
		mDeviceConfigPairingSecondaryInstruction: conf.ChipDeviceConfigPairingSecondaryInstruction,
		deviceConfigEnableCommissionableDeviceType: false,
	}
}

func (c ConfigurationManager) GetVendorId() (int, error) {
	return c.mVendorId, nil
}

func (c ConfigurationManager) GetVendorName() string {
	return c.mVendorName
}

func (c ConfigurationManager) GetProductId() (int, error) {
	return c.mProductId, nil
}

func (c ConfigurationManager) GetProductName() string {
	return c.mProductName
}

func (c ConfigurationManager) GetPrimaryMACAddress() (mac net.HardwareAddr, err error) {
	return c.GetPrimaryWiFiMACAddress(), nil
}

func (c ConfigurationManager) GetPrimaryWiFiMACAddress() (mac net.HardwareAddr) {
	ifs, _ := net.Interfaces()
	for _, i := range ifs {
		if !strings.Contains(i.Name, "lo") && i.HardwareAddr != nil {
			mac = i.HardwareAddr
		}
	}
	return
}

func (c ConfigurationManager) IsCommissionableDeviceTypeEnabled() bool {
	return c.deviceConfigEnableCommissionableDeviceType
}

func (c ConfigurationManager) GetDeviceTypeId() (device.MatterDeviceType, error) {
	return c.mDeviceType, nil
}

func (c ConfigurationManager) SetDeviceTypeId(t device.MatterDeviceType) {
	c.mDeviceType = t
}

func (c ConfigurationManager) IsCommissionableDeviceNameEnabled() bool {
	return true
}

func (c ConfigurationManager) GetCommissionableDeviceName() (string, error) {
	return c.mDeviceName, nil
}

func (c ConfigurationManager) GetInitialPairingHint() int {
	return c.mDevicePairingHint
}

func (c ConfigurationManager) GetInitialPairingInstruction() string {
	return c.mDeviceConfigPairingInitialInstruction
}

func (c ConfigurationManager) GetSecondaryPairingHint() int {
	return c.mDeviceSecondaryPairingHint
}

func (c ConfigurationManager) GetSecondaryPairingInstruction() string {
	return c.mDeviceConfigPairingSecondaryInstruction
}
