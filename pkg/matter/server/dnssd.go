package server

type Fabrics interface {
}

type DnssdServer struct {
	mSecuredPort   int
	mUnsecuredPort int
	mInterfaceId   any
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
	return nil
}
