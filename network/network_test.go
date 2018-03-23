package network_test

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"."
	"./mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type addr interface {
	Network() string
	String() string
}

func TestNewNetwork(t *testing.T) {
	networker := network.NewNetwork(nil, nil)
	assert.Implements(t, (*network.Network)(nil), networker, "NewNetwork() does not implement Network interface")
}

func TestStartHttpServerWithFileValidCase(t *testing.T) {
	httpServer := new(mocks.HttpServer)
	httpServer.On("Handle", mock.AnythingOfType("string"), mock.AnythingOfType("*http.fileHandler")).Return()
	httpServer.On("ListenAndServe", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	networker := network.NewNetwork(httpServer, nil)

	err := networker.StartHttpServer(8080, true)
	assert.NoError(t, err)
}

func TestStartHttpServerWithFileInvalidCase(t *testing.T) {
	httpServer := new(mocks.HttpServer)
	httpServer.On("Handle", mock.AnythingOfType("string"), mock.AnythingOfType("*http.fileHandler")).Return()
	httpServer.On("ListenAndServe", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("ListenAndServe Failed"))
	networker := network.NewNetwork(httpServer, nil)

	err := networker.StartHttpServer(8080, true)
	assert.Error(t, err, "StartHttpServer() not handling error properly")
}

func TestStartHttpServerWithoutFile(t *testing.T) {
	httpServer := new(mocks.HttpServer)
	httpServer.On("HandleFunc", mock.AnythingOfType("string"), mock.Anything).Return()
	httpServer.On("ListenAndServe", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	networker := network.NewNetwork(httpServer, nil)

	err := networker.StartHttpServer(8080, false)
	assert.NoError(t, err)
}

func TestAllUnavailablePortsWithEveryPortAvailable(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Close").Return(nil)
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	networker := network.NewNetwork(nil, localNet)

	list := networker.AllUnavailablePorts()
	assert.Empty(t, list)
}

func TestAllUnavailablePortsWithSomePortsUnavailable(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Close").Return(nil)
	localNet.On("Listen", mock.AnythingOfType("string"), ":3000").Return(nil, errors.New("Listen Failed"))
	localNet.On("Listen", mock.AnythingOfType("string"), ":4000").Return(nil, errors.New("Listen Failed"))
	localNet.On("Listen", mock.AnythingOfType("string"), ":5000").Return(nil, errors.New("Listen Failed"))
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	networker := network.NewNetwork(nil, localNet)

	list := networker.AllUnavailablePorts()
	assert.Equal(t, network.PortList{3000, 4000, 5000}, list)
}

func TestAllUnavailablePortsFromListWithEmptyInput(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Close").Return(nil)
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	networker := network.NewNetwork(nil, localNet)

	list := networker.AllUnavailablePortsFromList(&network.PortList{})
	assert.Empty(t, list)
}

func TestAllUnavailablePortsFromListWithEveryPortsAvailable(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Close").Return(nil)
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	networker := network.NewNetwork(nil, localNet)

	list := networker.AllUnavailablePortsFromList(&network.PortList{3000, 4000, 5000})
	assert.Empty(t, list)
}

func TestAllUnavailablePortsFromListWithSomePortsUnavailable(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Close").Return(nil)
	localNet.On("Listen", mock.AnythingOfType("string"), ":3000").Return(nil, errors.New("Listen Failed"))
	localNet.On("Listen", mock.AnythingOfType("string"), ":4000").Return(nil, errors.New("Listen Failed"))
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	networker := network.NewNetwork(nil, localNet)

	list := networker.AllUnavailablePortsFromList(&network.PortList{3000, 4000, 5000})
	assert.Equal(t, network.PortList{3000, 4000}, list)
}

func TestPortIsAvailableWithoutError(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Close").Return(nil)
	localNet.On("Listen", mock.AnythingOfType("string"), ":8080").Return(localNetListener, nil)
	networker := network.NewNetwork(nil, localNet)

	status, err := networker.PortIsAvailable(8080)
	assert.NoError(t, err)
	assert.True(t, status)
}

func TestPortIsAvailableWithError(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNet.On("Listen", mock.AnythingOfType("string"), ":8080").Return(nil, errors.New("Listen Failed"))
	networker := network.NewNetwork(nil, localNet)

	status, err := networker.PortIsAvailable(8080)
	assert.Error(t, err)
	assert.False(t, status)
}

func TestInternalIPWithoutError(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetConn := new(mocks.LocalNetConn)
	localNetConn.On("Close").Return(nil)
	localNetConn.On("LocalAddr").Return(&net.UDPAddr{})
	localNet.On("Dial", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetConn, nil)
	networker := network.NewNetwork(nil, localNet)
	_, err := networker.InternalIP()
	assert.NoError(t, err)
}

func TestInternalIPWithError(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNet.On("Dial", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Dial Failed"))
	networker := network.NewNetwork(nil, localNet)
	_, err := networker.InternalIP()
	assert.Error(t, err)
}

func TestExternalIPWithValidData(t *testing.T) {
	httpServer := new(mocks.HttpServer)
	httpServer.On("Get", mock.AnythingOfType("string")).Return(&http.Response{
		Body: nopCloser{bytes.NewBufferString("1.1.1.1")},
	}, nil)
	networker := network.NewNetwork(httpServer, nil)

	ip, err := networker.ExternalIP()
	assert.NoError(t, err)
	assert.NotNil(t, ip, "ExternalIP() failed to return IP with valid IP address")
}

func TestExternalIPWithInvalidData(t *testing.T) {
	httpServer := new(mocks.HttpServer)
	httpServer.On("Get", mock.AnythingOfType("string")).Return(&http.Response{
		Body: nopCloser{bytes.NewBufferString("invalid ip")},
	}, nil)
	networker := network.NewNetwork(httpServer, nil)

	ip, _ := networker.ExternalIP()
	assert.Nil(t, ip)
}

func TestExternalIPWithErrorResponse(t *testing.T) {
	httpServer := new(mocks.HttpServer)
	httpServer.On("Get", mock.AnythingOfType("string")).Return(nil, errors.New("Connection Failed"))
	networker := network.NewNetwork(httpServer, nil)

	ip, err := networker.ExternalIP()
	assert.Error(t, err, "ExternalIP() not handling error response from server properly")
	assert.Nil(t, ip)
}

func TestForwardingWithListenTCPFailed(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Listen Failed"))
	networker := network.NewNetwork(nil, localNet)

	err := networker.Forwarding("127.0.0.1:3000", 8080)
	assert.Error(t, err)
}

func TestForwardingWithListenerAcceptFailed(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Accept").Return(nil, errors.New("Accept Failed"))
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	networker := network.NewNetwork(nil, localNet)

	err := networker.Forwarding("127.0.0.1:3000", 8080)
	assert.Error(t, err)
}

func TestForwardingWithDialFailed(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetListener.On("Accept").Return(&net.TCPConn{}, nil)
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	localNet.On("Dial", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Dial Failed"))
	networker := network.NewNetwork(nil, localNet)

	err := networker.Forwarding("127.0.0.1:3000", 8080)
	assert.Error(t, err)
}

func TestForwardingWithNoError(t *testing.T) {
	localNet := new(mocks.LocalNet)
	localNetListener := new(mocks.LocalNetListener)
	localNetConn := new(mocks.LocalNetConn)
	localNetConn.On("Close").Return(nil)
	localNetListener.On("Accept").Return(&net.TCPConn{}, nil)
	localNet.On("Listen", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetListener, nil)
	localNet.On("Dial", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(localNetConn, nil)
	networker := network.NewNetwork(nil, localNet)

	errChan := make(chan error)
	go func() {
		errChan <- networker.Forwarding("127.0.0.1:3000", 8080)
	}()
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		return
	}
}
