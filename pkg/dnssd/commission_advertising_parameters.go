package dnssd

import (
	"github.com/galenliu/gateway/pkg/matter/core"
	"github.com/galenliu/gateway/pkg/matter/messageing"
	"net"
)

type Mac struct {
	mac string
}

type BaseAdvertisingParams struct {
	mPort         int
	mMac          string
	mEnableIPv4   bool
	mInterfaceId  net.Interface
	mMRPConfig    *messageing.ReliableMessageProtocolConfig
	mTcpSupported bool
}

type CommissionAdvertisingParameters struct {
	*BaseAdvertisingParams
	mVendorId          *uint16 //供应商口称
	mProductId         *uint16 //产品ID
	mDeviceType        *int32  //设备类型
	mPairingHint       int     //设备配提示
	mPairingInstr      string  //设备配对指南
	mDeviceName        string  //设备名称
	mMode              CommssionAdvertiseMode
	mCommissioningMode CommissioningMode
	mPeerId            *core.PeerId
}

type OperationalAdvertisingParameters struct {
	*BaseAdvertisingParams
	mPeerId core.PeerId
}

func (o *OperationalAdvertisingParameters) SetPeerId(peerId core.PeerId) {
	o.mPeerId = peerId
}

func (o *OperationalAdvertisingParameters) GetCompressedFabricId() core.CompressedFabricId {
	return o.mPeerId.GetCompressedFabricId()
}

func (o *OperationalAdvertisingParameters) GetPeerId() core.PeerId {
	return o.mPeerId
}

func (c *CommissionAdvertisingParameters) SetCommissioningMode(mode CommissioningMode) {
	c.mCommissioningMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissioningMode() CommissioningMode {
	return c.mCommissioningMode
}

func (c *CommissionAdvertisingParameters) SetCommissionAdvertiseMode(mode CommssionAdvertiseMode) {
	c.mMode = mode
}

func (c *CommissionAdvertisingParameters) GetCommissionAdvertiseMode() CommssionAdvertiseMode {
	return c.mMode
}

func (c *CommissionAdvertisingParameters) SetVendorId(id uint16) {
	c.mVendorId = &id
}

func (c *CommissionAdvertisingParameters) SetProductId(id uint16) {
	c.mProductId = &id
}

func (c *CommissionAdvertisingParameters) SetDeviceType(t int32) {
	c.mDeviceType = &t
}

func (c *CommissionAdvertisingParameters) SetDeviceName(name string) {
	c.mDeviceName = name
}

func (c *CommissionAdvertisingParameters) SetTcpSupported(b bool) {
	c.mTcpSupported = b
}

func (c *CommissionAdvertisingParameters) SetPairingHint(value int) {
	c.mPairingHint = value
}

func (c *CommissionAdvertisingParameters) SetPairingInstruction(ist string) {
	c.mPairingInstr = ist
}

func (c *CommissionAdvertisingParameters) SetMRPConfig(config *messageing.ReliableMessageProtocolConfig) {
	c.mMRPConfig = config
}

func (c *CommissionAdvertisingParameters) GetVendorId() *uint16 {
	return c.mVendorId
}

func (c *CommissionAdvertisingParameters) GetDeviceType() *int32 {
	return c.mDeviceType
}

func (c *CommissionAdvertisingParameters) GetProductId() *uint16 {
	return c.mProductId
}

func (c *CommissionAdvertisingParameters) GetDeviceName() string {
	return c.mDeviceName
}

func (b *BaseAdvertisingParams) IsIPv4Enabled() bool {
	return b.mEnableIPv4
}

func (b *BaseAdvertisingParams) SetPort(port int) {
	b.mPort = port
}

func (b *BaseAdvertisingParams) GetPort() int {
	return b.mPort
}

func (b *BaseAdvertisingParams) SetMaC(mac string) {
	b.mMac = mac
}

func (b *BaseAdvertisingParams) GetMac() string {
	if b.mMac == "" {
		b.mMac = mac48Address(randHex())
	}
	return b.mMac
}

func (b *BaseAdvertisingParams) GetUUID() string {
	if b.mMac == "" {
		b.mMac = mac48Address(randHex())
	}
	return b.mMac
}

func (b *BaseAdvertisingParams) EnableIpV4(enable bool) {
	b.mEnableIPv4 = enable
}
