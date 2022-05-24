package server

import "github.com/galenliu/gateway/pkg/dnssd"

type Fabrics interface {
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

func (d *DnssdServer) SetInterfaceId(id any) {
	d.mInterfaceId = id
}

func (d DnssdServer) StartServer() error {
	d.advertiser = dnssd.NewAdvertiser()

	err := d.advertiser.Init()
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
