package handlers

import (
	"io"
	"net"
	"net/url"
)

type AccessoriesHandler interface {
	HandleGetAccessories(i io.Reader) (io.Reader, error)
}

type CharacteristicHandler interface {
	HandleGetCharacteristics(url.Values, net.Conn) (io.Reader, error)
	HandleUpdateCharacteristics(io.Reader, net.Conn) error
}
