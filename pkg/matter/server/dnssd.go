package server

import (
	"github.com/galenliu/gateway/pkg/dnssd"
	"github.com/galenliu/gateway/pkg/matter/device"
)

type Fabrics interface {
	FabricCount() int
}

type DnssdServer struct {
	mSecuredPort   int
	mUnsecuredPort int
	mInterfaceId   any
	advertiser     *dnssd.Advertiser
	mFabrics       Fabrics
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

func (d *DnssdServer) SetInterfaceId(id any) {
	d.mInterfaceId = id
}

func (d DnssdServer) StartServer() error {
	d.advertiser = dnssd.NewAdvertiser()
	err := d.advertiser.Init(nil)
	if err != nil {
		return err
	}

	err = d.advertiser.RemoveServices()
	if err != nil {
		return err
	}

	err = d.AdvertiseOperational()
	if err != nil {
		return err
	}

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
