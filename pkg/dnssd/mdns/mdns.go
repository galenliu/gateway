package mdns

type InetLayer interface {
	NewUDPEndPoint()
}

type MdnsServer struct {
}

func NewMdnsServer() *MdnsServer {
	return &MdnsServer{}
}

func (m MdnsServer) Shutdown() {

}

func (m *MdnsServer) StartServer(inetLayer InetLayer, port int) error {
	m.Shutdown()
	return m.Listen(inetLayer, port)
}

func (m *MdnsServer) Listen(inetLayer InetLayer, port int) error {
	m.Shutdown()
	return nil
}
