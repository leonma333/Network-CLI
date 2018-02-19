package network

import (
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
	fakeNetwork.(*networkHandler).httpClient = &httpServerMock{}
	fakeNetwork.(*networkHandler).netClient = &localNetMock{}
}

/*
 * Turn error flags on for mock handlers
 */
func setupError() {
	fakeNetwork.(*networkHandler).httpClient.(*httpServerMock).err = true
	fakeNetwork.(*networkHandler).netClient.(*localNetMock).err = true
}

func TestNewNetwork(t *testing.T) {
	setupMock()
	if _, ok := fakeNetwork.(Network); !ok {
		t.Errorf("NewNetwork() does not implement Network interface")
	}
}

func TestStartHttpServerWithFile(t *testing.T) {
	setupMock()
	err := fakeNetwork.StartHttpServer(8080, true)
	if err != nil {
		t.Errorf("StartHttpServer() failed - %s", err.Error())
	}

	setupError()
	err = fakeNetwork.StartHttpServer(8080, true)
	if err == nil {
		t.Errorf("StartHttpServer() not handling error properly")
	}
}

func TestStartHttpServerWithoutFile(t *testing.T) {
	setupMock()
	err := fakeNetwork.StartHttpServer(8080, false)
	if err != nil {
		t.Errorf("StartHttpServer() failed - %s", err.Error())
	}

	setupError()
	err = fakeNetwork.StartHttpServer(8080, false)
	if err == nil {
		t.Errorf("StartHttpServer() not handling error properly")
	}
}

func TestAllUnavailablePorts(t *testing.T) {
	setupMock()
	list := fakeNetwork.AllUnavailablePorts()
	if !reflect.DeepEqual(list, mockNotAvailblePorts) {
		t.Errorf("AllUnavailablePorts() returns wrong value")
	}
}

func TestAllUnavailablePortsFromList(t *testing.T) {
	setupMock()
	testPorts := mockNotAvailblePorts
	testPorts = append(testPorts, 8080)
	list := fakeNetwork.AllUnavailablePortsFromList(&testPorts)
	if !reflect.DeepEqual(list, mockNotAvailblePorts) {
		t.Errorf("AllUnavailablePortsFromList() returns wrong value")
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
		t.Errorf("PortIsAvailable gives true for unavailable port")
	}

	setupError()
	_, err = fakeNetwork.PortIsAvailable(8080)
	if err == nil {
		t.Errorf("PortIsAvailable() not handling error properly")
	}
}
