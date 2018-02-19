package network

import (
	"errors"
	"net"
)

// localNet interface to define some local network functions
type localNet interface {
	Listen(protocol, address string) (net.Listener, error)
}

// localNetReal implements localNet with real network calls
type localNetReal struct{}

// localNetMock implements localNet with mock network calls
type localNetMock struct {
	err bool
}

/*
 * Real net.Listen from localNetReal (localNet)
 */
func (n *localNetReal) Listen(protocol, address string) (net.Listener, error) {
	return net.Listen(protocol, address)
}

/*
 * Mock net.Listen from httpServerMock (httpServer)
 */
func (n *localNetMock) Listen(protocol, address string) (net.Listener, error) {
	// if err flag is set, return error
	if n.err {
		return nil, errors.New("net.Listen error")
	}
	return nil, nil
}
