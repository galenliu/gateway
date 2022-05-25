package dnssd

type Mac struct {
	mac string
}

type BaseAdvertisingParams struct {
	mPort       int
	mMac        Mac
	mEnableIPv4 bool
}

type CommissionAdvertisingParameters struct {
	*BaseAdvertisingParams
	mVendorId          int    //供应商口称
	mProductId         int    //产品ID
	mDeviceType        int    //设备类型
	mPairingHint       int    //设备配提示
	mPairingInstr      string //设备配对指南
	mDeviceName        string //设备名称
	mMode              CommssionAdvertiseMode
	mCommissioningMode CommissioningMode
	mTcpSupported      bool
	mMac               Mac
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

func (c *CommissionAdvertisingParameters) SetVendorId(id int) {
	c.mVendorId = id
}

func (c *CommissionAdvertisingParameters) SetProductId(id int) {
	c.mProductId = id
}

func (c *CommissionAdvertisingParameters) SetDeviceType(t int) {
	c.mDeviceType = t
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

func (b *BaseAdvertisingParams) SetPort(port int) {
	b.mPort = port
}

func (b *BaseAdvertisingParams) GetPort() int {
	return b.mPort
}

func (b *BaseAdvertisingParams) SetMaC(mac string) {
	b.mMac = Mac{mac: mac}
}

func (b *BaseAdvertisingParams) GetMac() string {
	return b.mMac.mac
}

func (b *BaseAdvertisingParams) EnableIpV4(enable bool) {
	b.mEnableIPv4 = enable
}
