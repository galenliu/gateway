package server

import (
	"github.com/galenliu/gateway/pkg/dnssd"
	"github.com/galenliu/gateway/pkg/matter/config"
	"github.com/galenliu/gateway/pkg/matter/device"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/util"
	"net"
	"sync"
)

const Pkg = "Dnssd"

type Fabrics interface {
	FabricCount() int
}

type DnssdServer struct {
	mSecuredPort               int
	mUnsecuredPort             int
	mInterfaceId               net.Interface
	mCommissioningModeProvider *CommissioningWindowManager
	advertiser                 *dnssd.Advertiser
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

	d.advertiser = dnssd.NewAdvertiser()

	err := d.advertiser.Init(inet.UDPEndpointManager{})
	util.LogError(err, Pkg, "Failed initialize advertiser")

	err = d.advertiser.RemoveServices()
	util.LogError(err, Pkg, "Failed to remove advertised services")

	err = d.AdvertiseOperational()
	util.LogError(err, Pkg, "Failed to advertise operational node")

	if mode != dnssd.KDisabled {
		err = d.AdvertiseCommissionableNode(mode)
		util.LogError(err, Pkg, "Failed to advertise commissionable node")
	}

	// If any fabrics exist, the commissioning window must have been opened by the administrator
	// commissioning cluster commands which take care of the timeout.
	if !d.HaveOperationalCredentials() {
		d.ScheduleDiscoveryExpiration()
	}

	if config.ChipDeviceConfigEnableCommissionerDiscovery {
		err = d.AdvertiseCommissioner()
		util.LogError(err, Pkg, "Failed to advertise commissioner")
	}

	err = d.advertiser.FinalizeServiceUpdate()
	util.LogError(err, Pkg, "Failed to finalize service update")

	return nil
}

func (d DnssdServer) AdvertiseOperational() error {

	return nil
}

func (d DnssdServer) AdvertiseCommissioner() error {
	return nil
}

func (d DnssdServer) HaveOperationalCredentials() bool {
	return d.mFabrics.FabricCount() != 0
}

func (d DnssdServer) Advertise(commissionAbleNode bool, mode dnssd.CommissioningMode, config device.ConfigurationManager) error {
	advertiseParameters := dnssd.CommissionAdvertisingParameters{}
	if commissionAbleNode {
		advertiseParameters.SetPort(d.GetSecuredPort())
		advertiseParameters.SetCommissionAdvertiseMode(dnssd.CommissionableNode)
	} else {
		advertiseParameters.SetPort(d.GetUnsecuredPort())
		advertiseParameters.SetCommissionAdvertiseMode(dnssd.Commissioner)
	}
	advertiseParameters.SetCommissioningMode(mode)

	advertiseParameters.EnableIpV4(true)
	advertiseParameters.SetVendorId(config.GetVendorId())
	advertiseParameters.SetProductId(config.GetProductId())

	if config.IsCommissionableDeviceTypeEnabled() {
		advertiseParameters.SetDeviceType(config.GetDeviceTypeId())
	}
	if config.IsCommissionableDeviceNameEnabled() {
		name := config.GetCommissionableDeviceName()
		advertiseParameters.SetDeviceName(name)
	}
	advertiseParameters.SetTcpSupported(true)

	if !d.HaveOperationalCredentials() {
		value := config.GetInitialPairingHint()
		advertiseParameters.SetPairingHint(value)

		ist := config.GetInitialPairingInstruction()
		advertiseParameters.SetPairingInstruction(ist)
	} else {
		hint := config.GetSecondaryPairingHint()
		advertiseParameters.SetPairingHint(hint)

		ins := config.GetSecondaryPairingInstruction()
		advertiseParameters.SetPairingInstruction(ins)
	}
	mdnsAdvertiser := dnssd.NewAdvertiser()
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
