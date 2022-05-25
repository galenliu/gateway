package device

import (
	"net"
	"strings"
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

type ConfigurationManager struct {
	mVendorId                                  int
	mVendorName                                string
	mProductName                               string
	mProductId                                 int
	mDeviceType                                MatterDeviceType
	mDeviceName                                string
	mTcpSupported                              bool
	mDevicePairingHint                         int
	mDevicePairingSecondaryHint                int
	mDeviceSecondaryPairingHint                int
	mDeviceConfigPairingInitialInstruction     string
	mDeviceConfigPairingSecondaryInstruction   string
	deviceConfigEnableCommissionableDeviceType bool
}

func NewConfigurationManager(conf Config) ConfigurationManager {
	return ConfigurationManager{
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

func (c ConfigurationManager) GetVendorId() int {
	return c.mVendorId
}

func (c ConfigurationManager) GetVendorName() string {
	return c.mVendorName
}

func (c ConfigurationManager) GetProductId() int {
	return c.mProductId
}

func (c ConfigurationManager) GetProductName() string {
	return c.mProductName
}

func (c ConfigurationManager) GetPrimaryMACAddress() (mac net.HardwareAddr) {
	return c.GetPrimaryWiFiMACAddress()
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

func (c ConfigurationManager) GetDeviceTypeId() MatterDeviceType {
	return c.mDeviceType
}

func (c ConfigurationManager) SetDeviceTypeId(t MatterDeviceType) {
	c.mDeviceType = t
}

func (c ConfigurationManager) IsCommissionableDeviceNameEnabled() bool {
	return true
}

func (c ConfigurationManager) GetCommissionableDeviceName() string {
	return c.mDeviceName
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
