package network

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

var mockNotAvailblePorts = PortList([]int{1000, 2000, 3000})

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

	// if address is not available, return error
	for _, port := range mockNotAvailblePorts {
		if address == ":"+strconv.Itoa(port) {
			return nil, errors.New(fmt.Sprintf("net.Listen error: %s not available", address))
		}
	}

	return nil, nil
}
