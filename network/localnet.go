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
	Listen(protocol, address string) (localNetListener, error)
	// Dial(protocol, address string) (net.Conn, error)
}

// localNetListener interface to defince local network listener functions
type localNetListener interface {
	Accept() (net.Conn, error)
	Close() error
	Addr() net.Addr
}

// localNetReal implements localNet with real network calls
type localNetReal struct{}

// localNetMock implements localNet with mock network calls
type localNetMock struct {
	err bool
}

// localNetListenerMock implements localNetListener with mock network calls
type localNetListenerMock struct{}

/*
 * Real net.Listen from localNetReal (localNet)
 */
func (n *localNetReal) Listen(protocol, address string) (localNetListener, error) {
	return net.Listen(protocol, address)
}

/*
 * Real net.Dial from localNetReal (localNet)
 */
// func (n *localNetReal) Dial(protocol, address string) (net.Conn, error) {
// 	return net.Dial(protocol, address)
// }

/*
 * Mock net.Listen from httpServerMock (httpServer)
 */
func (n *localNetMock) Listen(protocol, address string) (localNetListener, error) {
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

	return localNetListenerMock{}, nil
}

/*
 * Mock net.Dial from localNetMock (localNet)
 */
// func (n *localNetMock) Dial(protocol, address string) (net.Conn, error) {
// 	return nil, nil
// }

/*
 * Mock net.Listener.Accept from localNetListenerMock (localNetListener)
 */
func (localNetListenerMock) Accept() (net.Conn, error) {
	return nil, nil
}

/*
 * Mock net.Listener.Close from localNetListenerMock (localNetListener)
 */
func (localNetListenerMock) Close() error {
	return nil
}

/*
 * Mock net.Listener.Addr from localNetListenerMock (localNetListener)
 */
func (localNetListenerMock) Addr() net.Addr {
	return nil
}
