package server

import (
	"github.com/galenliu/gateway/pkg/dnssd"
	"github.com/galenliu/gateway/pkg/matter/config"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/matter/messageing"
	"github.com/galenliu/gateway/pkg/matter/platform"
	"github.com/galenliu/gateway/pkg/util"
	"net"
	"sync"
)

type Fabrics interface {
	FabricCount() int
}

type DnssdServer struct {
	mSecuredPort               int
	mUnsecuredPort             int
	mInterfaceId               net.Interface
	mCommissioningModeProvider *CommissioningWindowManager
	mFabrics                   Fabrics
}

var insDnssd *DnssdServer
var onceDnssd sync.Once

func DnssdInstance() *DnssdServer {
	onceDnssd.Do(func() {
		insDnssd = NewDnssdServer()
	})
	return insDnssd
}

func NewDnssdServer() *DnssdServer {
	return &DnssdServer{}
}

func (d DnssdServer) Shutdown() {

}

func (d DnssdServer) SetFabricTable(f Fabrics) {
	d.mFabrics = f
}

func (d *DnssdServer) SetSecuredPort(port int) {
	d.mSecuredPort = port
}

func (d *DnssdServer) SetUnsecuredPort(port int) {
	d.mUnsecuredPort = port
}

func (d *DnssdServer) GetSecuredPort() int {
	return d.mSecuredPort
}

func (d *DnssdServer) GetUnsecuredPort() int {
	return d.mUnsecuredPort
}

func (d *DnssdServer) SetInterfaceId(inter net.Interface) {
	d.mInterfaceId = inter
}

func (d *DnssdServer) StartServer() error {
	var mode = dnssd.KDisabled
	if d.mCommissioningModeProvider != nil {
		mode = d.mCommissioningModeProvider.GetCommissioningMode()
	}
	return d.startServer(mode)
}

func (d *DnssdServer) startServer(mode dnssd.CommissioningMode) error {

	//使用UDPEndPointManager初始化一个Dnssd-Advertiser
	err := dnssd.AdvertiserInstance().Init(inet.UDPEndpoint{})
	util.LogError(err, "Discover", "Failed initialize advertiser")

	err = dnssd.AdvertiserInstance().RemoveServices()
	util.LogError(err, "Discover", "Failed to remove advertised services")

	err = d.AdvertiseOperational()
	util.LogError(err, "Discover", "Failed to advertise operational node")

	if mode != dnssd.KDisabled {
		err = d.AdvertiseCommissionableNode(mode)
		util.LogError(err, "Discover", "Failed to advertise commissionable node")
	}

	// If any fabrics exist, the commissioning window must have been opened by the administrator
	// commissioning cluster commands which take care of the timeout.
	if !d.HaveOperationalCredentials() {
		d.ScheduleDiscoveryExpiration()
	}

	if config.ChipDeviceConfigEnableCommissionerDiscovery {
		err = d.AdvertiseCommissioner()
		util.LogError(err, "Discover", "Failed to advertise commissioner")
	}

	err = dnssd.AdvertiserInstance().FinalizeServiceUpdate()
	util.LogError(err, "Discover", "Failed to finalize service update")

	return nil
}

func (d DnssdServer) AdvertiseOperational() error {

	return nil
}

func (d DnssdServer) AdvertiseCommissioner() error {
	return d.Advertise(false, dnssd.KDisabled)
}

func (d DnssdServer) HaveOperationalCredentials() bool {
	return d.mFabrics.FabricCount() != 0
}

func (d DnssdServer) Advertise(commissionAbleNode bool, mode dnssd.CommissioningMode) error {

	advertiseParameters := dnssd.CommissionAdvertisingParameters{}

	advertiseParameters.SetPort(util.ConditionFunc(commissionAbleNode, d.GetUnsecuredPort, d.GetUnsecuredPort))
	advertiseParameters.SetCommissionAdvertiseMode(util.ConditionValue(commissionAbleNode, dnssd.KCommissionableNode, dnssd.KCommissioner))

	advertiseParameters.SetCommissioningMode(mode)

	mac, err := platform.ConfigurationMgr().GetPrimaryMACAddress()
	util.LogError(err, "Discovery", "Failed to get primary mac address of device. Generating a random one.")
	advertiseParameters.SetMaC(util.ConditionValue(err != nil, mac, util.GenerateMac()))

	vid, err := platform.ConfigurationMgr().GetVendorId()
	util.LogError(err, "Discovery", "Vendor ID not known")
	if err != nil {
		advertiseParameters.SetVendorId(vid)
	}

	pid, err := platform.ConfigurationMgr().GetProductId()
	util.LogError(err, "Discovery", "Product ID not known")
	if err != nil {
		advertiseParameters.SetProductId(pid)
	}

	//uint16_t discriminator = 0;
	//CHIP_ERROR error       = DeviceLayer::GetCommissionableDataProvider()->GetSetupDiscriminator(discriminator);
	//if (error != CHIP_NO_ERROR)
	//{
	//	ChipLogError(Discovery,
	//		"Setup discriminator read error (%" CHIP_ERROR_FORMAT ")! Critical error, will not be commissionable.",
	//	error.Format());
	//	return error;
	//}

	// Override discriminator with temporary one if one is set
	//discriminator = mEphemeralDiscriminator.ValueOr(discriminator);
	//
	//advertiseParameters.SetShortDiscriminator(static_cast<uint8_t>((discriminator >> 8) & 0x0F))
	//.SetLongDiscriminator(discriminator);
	//

	if platform.ConfigurationMgr().IsCommissionableDeviceTypeEnabled() {
		did, err := platform.ConfigurationMgr().GetDeviceTypeId()
		if err != nil {
			advertiseParameters.SetDeviceType(did)
		}
	}

	if platform.ConfigurationMgr().IsCommissionableDeviceNameEnabled() {
		name, err := platform.ConfigurationMgr().GetCommissionableDeviceName()
		if err != nil {
			advertiseParameters.SetDeviceName(name)
		}
	}

	advertiseParameters.SetMRPConfig(messageing.GetLocalMRPConfig())
	advertiseParameters.SetTcpSupported(inet.InetConfigEnableTcpEndpoint)

	if !d.HaveOperationalCredentials() {
		value := platform.ConfigurationMgr().GetInitialPairingHint()

		advertiseParameters.SetPairingHint(value)

		ist := platform.ConfigurationMgr().GetInitialPairingInstruction()
		advertiseParameters.SetPairingInstruction(ist)
	} else {
		hint := platform.ConfigurationMgr().GetSecondaryPairingHint()
		advertiseParameters.SetPairingHint(hint)

		ins := platform.ConfigurationMgr().GetSecondaryPairingInstruction()
		advertiseParameters.SetPairingInstruction(ins)
	}
	mdnsAdvertiser := dnssd.AdvertiserInstance()
	return mdnsAdvertiser.Advertise(advertiseParameters)
}

func (d *DnssdServer) SetCommissioningModeProvider(manager *CommissioningWindowManager) {
	d.mCommissioningModeProvider = manager
}

func (d *DnssdServer) AdvertiseCommissionableNode(mode dnssd.CommissioningMode) error {
	return nil
}

func (d DnssdServer) ScheduleDiscoveryExpiration() {
	//TODO
	return
}
