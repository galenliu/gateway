package udp_endpoint

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/system"
	"net"
	"net/netip"
)

type State uint8

type UDPEndpointImpl interface {
	Close()
	Bind(addr netip.Addr, port int, interfaceId inet.InterfaceId) error
	Listen(funct OnMessageReceivedFunct, errorFunct OnReceiveErrorFunct, appState any) error
	SendTo(addr netip.Addr, port int, handle *system.PacketBufferHandle, interfaceId inet.InterfaceId) error
	SendMsg(pktInfo *inet.IPPacketInfo, msg *system.PacketBufferHandle) error
	LeaveMulticastGroup(interfaceId inet.InterfaceId, addr netip.Addr) error
}

type OnMessageReceivedFunct = func(*system.PacketBufferHandle, *inet.IPPacketInfo)
type OnReceiveErrorFunct = func(error, *inet.IPPacketInfo)

type impl interface {
	BindImpl(addr netip.Addr, port int, interfaceId inet.InterfaceId) error
	SendMsgImpl(pktInfo *inet.IPPacketInfo, msg *system.PacketBufferHandle) error
	IPv4JoinLeaveMulticastGroupImpl(aInterfaceId inet.InterfaceId, addr netip.Addr, b bool) error
	IPv6JoinLeaveMulticastGroupImpl(aInterfaceId inet.InterfaceId, addr netip.Addr, b bool) error
	ListenImpl() error
	CloseImpl()
}

const (
	kReady State = iota
	kBound
	kListening
	kClosed
)

type UDPEndpoint struct {
	mInterface net.Interface
	mAddr      netip.Addr
	mPort      int
	mState     State
	mAppState  any
	impl
	onMessageReceived OnMessageReceivedFunct
	onReceiveError    OnReceiveErrorFunct
}

// DefaultUDPEndpoint 初始化一个默认Socket的UDPEndPoint
func DefaultUDPEndpoint() *UDPEndpoint {
	up := &UDPEndpoint{}
	up.impl = NewUDPEndPointImplSockets()
	return up
}

func (e *UDPEndpoint) Bind(addr netip.Addr, port int, interfaceId inet.InterfaceId) error {
	if e.mState != kReady && e.mState != kBound {
		return errors.IncorrectState("not ready or bound")
	}
	if e.impl == nil {
		return errors.NotImplement("UDPEndpoint")
	}
	err := e.BindImpl(addr, port, interfaceId)
	if err != nil {
		return err
	}
	e.mState = kBound
	return nil
}

func (e *UDPEndpoint) Listen(funct OnMessageReceivedFunct, errorFunct OnReceiveErrorFunct, appState any) error {
	if e.mState == kListening {
		return nil
	}
	if e.mState != kBound {
		return errors.IncorrectState("not bound")
	}
	e.onMessageReceived = funct
	e.onReceiveError = errorFunct
	e.mAppState = appState
	err := e.ListenImpl()
	if err != nil {
		return err
	}
	return nil
}

func (e *UDPEndpoint) SendTo(addr netip.Addr, port int, msg *system.PacketBufferHandle, interfaceId inet.InterfaceId) error {
	pktInfo := &inet.IPPacketInfo{
		DestAddress: addr,
		InterfaceId: interfaceId,
		DestPort:    port,
	}
	return e.SendMsg(pktInfo, msg)
}

func (e *UDPEndpoint) SendMsg(pktInfo *inet.IPPacketInfo, msg *system.PacketBufferHandle) error {
	return e.SendMsgImpl(pktInfo, msg)
}

func (e *UDPEndpoint) JoinMulticastGroup(interfaceId inet.InterfaceId, addr netip.Addr) error {
	if !addr.IsMulticast() {
		return fmt.Errorf("wrong address type")
	}
	if addr.Is4() {
		return e.IPv4JoinLeaveMulticastGroupImpl(interfaceId, addr, true)
	}
	if addr.Is6() {
		return e.IPv6JoinLeaveMulticastGroupImpl(interfaceId, addr, true)
	}
	return fmt.Errorf("wrong address type")
}

func (e *UDPEndpoint) LeaveMulticastGroup(interfaceId inet.InterfaceId, addr netip.Addr) error {
	if !addr.IsMulticast() {
		return fmt.Errorf("wrong address type")
	}
	if addr.Is4() {
		return e.IPv4JoinLeaveMulticastGroupImpl(interfaceId, addr, false)
	}
	if addr.Is6() {
		return e.IPv6JoinLeaveMulticastGroupImpl(interfaceId, addr, false)
	}

	return fmt.Errorf("wrong address type")
}

func (e *UDPEndpoint) Close() {
	if e.mState != kClosed {
		e.CloseImpl()
		e.mState = kClosed
	}
}
