package main

import (
	"context"
	"fmt"
	"github.com/brutella/dnssd"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/xiam/to"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func DnssdStart() error {

	// Listen with a tcp socket on a given addr/port.
	tcpLn, err := net.Listen("tcp", "")
	if err != nil {
		return err
	}
	log.Print("addr:" + tcpLn.Addr().String() + "\t\n")

	// Get the port from the listener address because it
	// it might be different than specified in Port.
	_, port, _ := net.SplitHostPort(tcpLn.Addr().String())

	i, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	cfg := dnssd.Config{
		Name:   "Lamp123",
		Type:   "_hap._tcp",
		Domain: "local",
		Host:   strings.Replace(mac48Address(randHex()), ":", "", -1), // use the id (without the colons) to get unique hostnames
		Text:   txtRecords(),
		Port:   i,
	}

	rp, _ := dnssd.NewResponder()

	sv, _ := dnssd.NewService(cfg)

	_, _ = rp.Add(sv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dnsCtx, dnsCancel := context.WithCancel(ctx)
	defer dnsCancel()

	dnsStop := make(chan struct{})
	//go func() {
	rp.Respond(dnsCtx)
	log.Print("dnssd responder stopped")
	dnsStop <- struct{}{}
	//}()
	return nil
}

func txtRecords() map[string]string {
	return map[string]string{
		"pv": "1.0",
		"id": "D0:5A:97:7F:34:80",
		"c#": fmt.Sprintf("%d", 2),
		"s#": "1",
		"sf": fmt.Sprintf("%d", to.Int64(true)),  //isPaired
		"ff": fmt.Sprintf("%d", to.Int64(false)), //MfiCompliant
		"md": "lamp",
		"ci": fmt.Sprintf("%d", 8), //Info Type
		"sh": "F2YhbQ==",
	}
}

func HomeKitStart() {
	a := accessory.NewSwitch(accessory.Info{
		Name: "Lamp",
	})

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	// Create the hap server.
	server, err := hap.NewServer(fs, a.A)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Run the server.
	server.ListenAndServe(ctx)
}

func main() {
	r := randHex()
	log.Print(r + "\t\n")
	log.Print(mac48Address(r))
	DnssdStart()
	//HomeKitStart()
}
