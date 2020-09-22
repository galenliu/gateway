package hap

import "net"

type Session interface {
	Connection() net.Conn
}

type session struct {
	connection net.Conn
}

func NewSession(conn net.Conn) Session {
	s := session{connection: conn}
	return &s

}

func (s *session) Connection() net.Conn {
	return s.connection
}
