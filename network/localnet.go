package network

import (
	"net"
	"time"
)

// LocalNet interface to define some local network functions
type LocalNet interface {
	Listen(protocol, address string) (LocalNetListener, error)
	Dial(protocol, address string) (LocalNetConn, error)
}

// LocalNetListener interface to defince local network listener functions
type LocalNetListener interface {
	Accept() (net.Conn, error)
	Close() error
	Addr() net.Addr
}

// LocalNetConn interface to defince local network connection functions
type LocalNetConn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type localNet struct{}

func (n *localNet) Listen(protocol, address string) (LocalNetListener, error) {
	return net.Listen(protocol, address)
}

func (n *localNet) Dial(protocol, address string) (LocalNetConn, error) {
	return net.Dial(protocol, address)
}
