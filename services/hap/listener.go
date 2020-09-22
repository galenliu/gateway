package hap

import (
	"net"
)

type Listener struct {
	*net.TCPListener
	context Context
}

func NewHapListener(l *net.TCPListener, ctx Context) *Listener {
	return &Listener{l, ctx}
}

func (l *Listener) Accept() (c net.Conn, err error) {

	conn, err := l.AcceptTCP()
	if err != nil {
		return
	}
	hapConn := NewConnection(conn, l.context)
	return hapConn, err
}
