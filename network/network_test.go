package network

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	fakeNetwork = NewNetwork()
)

/*
 * Setup mock handlers
 */
func setupMock() {
	fakeNetwork.(*networkHandler).httpClient = &httpServerMock{
		body: readerMock{bytes.NewBufferString("69.89.31.226\n")},
	}
	fakeNetwork.(*networkHandler).netClient = &localNetMock{}
}

/*
 * Turn error flags on for mock handlers
 */
func setupError(wrongOutput bool) {
	if wrongOutput {
		fakeNetwork.(*networkHandler).httpClient.(*httpServerMock).body = readerMock{bytes.NewBufferString("invalid_ip\n")}
	} else {
		fakeNetwork.(*networkHandler).httpClient.(*httpServerMock).err = true
	}
	fakeNetwork.(*networkHandler).netClient.(*localNetMock).err = true
}

func TestNewNetwork(t *testing.T) {
	if _, ok := fakeNetwork.(Network); !ok {
		t.Error("NewNetwork() does not implement Network interface")
	}
}

func TestStartHttpServerWithFile(t *testing.T) {
	setupMock()
	err := fakeNetwork.StartHttpServer(8080, true)
	if err != nil {
		t.Errorf("StartHttpServer() failed - %s", err.Error())
	}

	setupError(false)
	err = fakeNetwork.StartHttpServer(8080, true)
	if err == nil {
		t.Error("StartHttpServer() not handling error properly")
	}
}

func TestStartHttpServerWithoutFile(t *testing.T) {
	setupMock()
	err := fakeNetwork.StartHttpServer(8080, false)
	if err != nil {
		t.Errorf("StartHttpServer() failed - %s", err.Error())
	}

	setupError(false)
	err = fakeNetwork.StartHttpServer(8080, false)
	if err == nil {
		t.Error("StartHttpServer() not handling error properly")
	}
}

func TestAllUnavailablePorts(t *testing.T) {
	setupMock()
	list := fakeNetwork.AllUnavailablePorts()
	if !reflect.DeepEqual(list, mockNotAvailblePorts) {
		t.Error("AllUnavailablePorts() returns wrong value")
	}
}

func TestAllUnavailablePortsFromList(t *testing.T) {
	setupMock()
	testPorts := mockNotAvailblePorts
	testPorts = append(testPorts, 8080)
	list := fakeNetwork.AllUnavailablePortsFromList(&testPorts)
	if !reflect.DeepEqual(list, mockNotAvailblePorts) {
		t.Error("AllUnavailablePortsFromList() returns wrong value")
	}
}

func TestPortIsAvailable(t *testing.T) {
	setupMock()
	status, err := fakeNetwork.PortIsAvailable(8080)
	if !status || err != nil {
		t.Errorf("PortIsAvailable() failed - %s", err.Error())
	}

	status, _ = fakeNetwork.PortIsAvailable(mockNotAvailblePorts[0])
	if status {
		t.Error("PortIsAvailable gives true for unavailable port")
	}

	setupError(false)
	_, err = fakeNetwork.PortIsAvailable(8080)
	if err == nil {
		t.Error("PortIsAvailable() not handling error properly")
	}
}

func TestInternalIP(t *testing.T) {
	_, err := fakeNetwork.InternalIP()
	if err != nil {
		t.Errorf("InternalIP() failed - %s", err.Error())
	}
}

func TestExternalIPWithValidData(t *testing.T) {
	setupMock()
	ip, err := fakeNetwork.ExternalIP()
	if err != nil {
		t.Errorf("ExternalIP() failed - %s", err.Error())
	}
	if ip == nil {
		t.Error("ExternalIP() failed to return IP with valid IP address")
	}
}

func TestExternalIPWithInvalidData(t *testing.T) {
	setupMock()
	setupError(true)
	ip, _ := fakeNetwork.ExternalIP()
	if ip != nil {
		t.Error("ExternalIP() not handling invalid IP properly")
	}
}

func TestExternalIPWithErrorResponse(t *testing.T) {
	setupMock()
	setupError(false)
	_, err := fakeNetwork.ExternalIP()
	if err == nil {
		t.Error("ExternalIP() not handling error response properly")
	}
}
