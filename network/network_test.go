package network

import "testing"

var (
	fakeNetwork = NewNetwork()
)

/*
 * Setup mock handlers
 */
func setupMock() {
	fakeNetwork.(*networkHandler).httpClient = &httpServerMock{}
}

/*
 * Turn error flags on for mock handlers
 */
func setupError() {
	fakeNetwork.(*networkHandler).httpClient.(*httpServerMock).err = true
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
		t.Errorf("StartHttpServer failed - %s", err.Error())
	}

	setupError()
	err = fakeNetwork.StartHttpServer(8080, false)
	if err == nil {
		t.Errorf("StartHttpServer() not handling error properly")
	}
}
