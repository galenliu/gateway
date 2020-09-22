package zeroconf

import (
	"context"
	"fmt"
	"github.com/grandcat/zeroconf"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DnsSdConf struct {
	Name   string //service name
	Type   string //type is the service of type,for  example _hap._tcp
	Domain string //domain is the name of domain, for example "local", if empty "local" used
	Host   string

	Text map[string]string
	IPs  []net.IP
	Port int //Port is the port of the service

	ifaceIPs map[string][]net.IP
}

type DnsSd struct {
	Name     string
	Type     string
	Domain   string
	Host     string
	Text     map[string]string
	TTL      time.Duration
	Port     int
	IPs      []net.IP
	IfaceIPs map[string][]net.IP

	expiration time.Time
}

func (zc *DnsSd) DiscoverServices() {

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println(entry)
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, "_workstation._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}
	<-ctx.Done()
}

func (zc *DnsSd) RegisterService(c *config) {

	server, err := zeroconf.Register(
		"GoZeroconf",
		"_workstation._tcp",
		"local.",
		42424,
		[]string{"txtv=0", "lo=1", "la=2"},
		nil)

	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-zc.flag:
		print("exit by chan")
	case <-sig:
		// Exit by user
	case <-time.After(time.Second * 120):
		print("exit by time out")
	}

	log.Println("Shutting down.")
}

func (zc *DnsSd) Stop() {

}

func NewServer(c *DnsSdConf) (s Service, err error) {
	name := c.Name
	typ := c.Type
	port := c.Port

	if len(name) == 0 {
		err = fmt.Errorf("invaild name \"%s\"", name)
		return
	}

	if len(typ) == 0 {
		err = fmt.Errorf("invaild type \"%s\"", typ)
		return
	}

	if port == 0 {
		err = fmt.Errorf("invaild name \"%d\"", port)
		return
	}

	domain := c.Domain
	if len(domain) == 0 {
		domain = "local"
	}

	host := c.Host
	if len(host) == 0 {
		host = hostname()
	}

	text := c.Text
	if text == nil {
		text = map[string]string{}
	}

	ips := []net.IP{}
	ifaceIPs := map[string][]net.IP{}

	if c.IPs != nil && len(c.IPs) > 0 {
		ips = c.IPs
	}

	if c.ifaceIPs != nil && len(c.ifaceIPs) > 0 {
		ifaceIPs = c.ifaceIPs
	}

	return &DnsSd{
		Name:     name,
		Type:     typ,
		Domain:   domain,
		Host:     host,
		Text:     text,
		IPs:      ips,
		IfaceIPs: ifaceIPs,
	}, err
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}
