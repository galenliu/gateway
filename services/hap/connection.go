package hap

import (
	"io"
	"net"
	"time"
)

type Connection struct {
	connection net.Conn
	context    Context

	readBuffer io.Reader
}

func NewConnection(connection net.Conn, ctx Context) *Connection {
	conn := &Connection{connection: connection, context: ctx}
	session := NewSession(conn)
	ctx.SetSessionForConnection(session, conn)
	return conn
}

func (c *Connection) Read(b []byte) (n int, err error) {
	return c.connection.Read(b)
}

func (c *Connection) Write(b []byte) (n int, err error) {
	return c.connection.Write(b)
}

func (c *Connection) Close() error {
	return c.connection.Close()
}

func (c *Connection) LocalAddr() net.Addr {
	return c.connection.LocalAddr()
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.connection.RemoteAddr()
}

func (c *Connection) SetDeadline(t time.Time) error {
	return c.connection.SetDeadline(t)
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.connection.SetDeadline(t)
}

func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.connection.SetWriteDeadline(t)
}
