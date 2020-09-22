package http

import (
	"gateway/accessory"
	"gateway/services/hap"
	"gateway/services/hap/controller"
	"gateway/services/hap/handlers"
	"log"
	"net"
	"net/http"
	"sync"
)

type Config struct {
	Port      string
	Mutex     *sync.Mutex
	Context   hap.Context
	Container *accessory.Container
}

type Server struct {
	port      string
	mutex     *sync.Mutex
	context   hap.Context
	container *accessory.Container
	listener  *net.TCPListener

	hapListener *hap.Listener

	Mux *http.ServeMux
}

func NewServer(c Config) *Server {

	ln, err := net.Listen("tcp", c.Port)
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	s := Server{
		port:      port,
		mutex:     c.Mutex,
		context:   c.Context,
		container: c.Container,
		listener:  ln.(*net.TCPListener),
	}
	return &s
}

func (s *Server) ListenAndServe() {

	s.listenAndServe(s.addString())
	s.setHandle()
}

func (s *Server) listenAndServe(add string) error {
	var server = http.Server{Addr: add, Handler: s.Mux}
	var listener = hap.NewHapListener(s.listener, s.context)
	s.hapListener = listener
	return server.Serve(s.hapListener)

}

func (s *Server) addString() string {
	return ":" + s.port
}

func (s *Server) setHandle() {
	accessoriesController := controller.NewAccessoriesController(s.container)
	characteristicController := controller.NewCharacteristicController(s.container)

	s.Mux.Handle("/pair-setup", handlers.NewPairSetup())
	s.Mux.Handle("/accessories", handlers.NewAccessories(s.mutex, accessoriesController))
	s.Mux.Handle("/characteristics", handlers.NewCharacteristics(s.mutex, characteristicController, s.context))
}
