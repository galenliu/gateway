package dnssd

import "github.com/galenliu/gateway/pkg/dnssd/mdns"

const MdnsPort = 5353

type MdnsServerBase interface {
	Shutdown()
}

type Advertiser struct {
	mResponseSender *ResponseSender
	mdnsServer      *mdns.MdnsServer
}

func NewAdvertiser() *Advertiser {
	return &Advertiser{}
}

func (s *Advertiser) Init(layer mdns.InetLayer) error {

	s.mdnsServer = mdns.NewMdnsServer()
	s.mdnsServer.Shutdown()

	s.mResponseSender = NewResponseSender()
	s.mResponseSender.SetServer(s.mdnsServer)

	err := s.mdnsServer.StartServer(layer, MdnsPort)
	if err != nil {
		return err
	}

	s.AdvertiseRecords()

	return nil
}

func (s *Advertiser) RemoveServices() error {
	return nil
}

func (s *Advertiser) AdvertiseRecords() {

}
