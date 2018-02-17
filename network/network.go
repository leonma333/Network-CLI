package network

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

const (
	DEFAULT_PORT = 8080
	MAX_PORT     = 65535
)

// Network interface to deal with http server and IP
type Network interface {
	StartHttpServer(port int, useFile bool) error
	AllUnavailablePorts() PortList
	AllUnavailablePortsFromList(pl *PortList) PortList
	PortIsAvailable(port int) (status bool, err error)
	InternalIP() (net.IP, error)
	ExternalIP() (net.IP, error)
	Forwarding(target string, port int) error
	forward(conn net.Conn, target string)
}

// NetworkHelper implements Network interface
type networkHandler struct{}

func NewNetwork() Network {
	return &networkHandler{}
}

/*
 * Start a local http server with the port number provided
 */
func (network *networkHandler) StartHttpServer(port int, useFile bool) error {
	if useFile {
		http.Handle("/", http.FileServer(http.Dir("./")))
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "You're now on port %v [%s]", port, r.URL.Path[0:])
		})
	}

	//go func() {
	if err := http.ListenAndServe(portStringify(port), nil); err != nil {
		return err
	}
	//}()
	return nil
}

/*
 * Return all the unavailable port number from the machine
 */
func (network *networkHandler) AllUnavailablePorts() PortList {
	var unavailablePorts PortList
	for i := 0; i <= MAX_PORT; i++ {
		if status, _ := network.PortIsAvailable(i); !status {
			unavailablePorts = append(unavailablePorts, i)
		}
	}
	return unavailablePorts
}

/*
 * Return all the unavailable port number from the given list of port numbers
 */
func (network *networkHandler) AllUnavailablePortsFromList(pl *PortList) PortList {
	var unavailablePorts PortList
	for _, port := range *pl {
		if status, _ := network.PortIsAvailable(port); !status {
			unavailablePorts = append(unavailablePorts, port)
		}
	}
	return unavailablePorts
}

/*
 * Check the given port number is availble to be used for the machine
 */
func (network *networkHandler) PortIsAvailable(port int) (status bool, err error) {
	host := ":" + strconv.Itoa(port)
	server, err := net.Listen("tcp", host)
	if err != nil {
		return false, err
	}
	server.Close()
	return true, nil
}

/*
 * Retrieve machine local IP address
 */
func (network *networkHandler) InternalIP() (net.IP, error) {
	// Dial to connect to local server
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}

	// Close the connection when the response is read
	defer conn.Close()

	// Get the local IP address
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

/*
 * Retrieve machine external IP address
 */
func (network *networkHandler) ExternalIP() (net.IP, error) {
	// Get request on myexternalip.com to retrieve external IP address for the machine
	resp, err := http.Get("http://myexternalip.com/raw")

	if err != nil {
		return nil, err
	}

	// Close the connection when response body is read
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	ipStr := string(body)

	// Convert body byte array to string and return it
	return net.ParseIP(ipStr[:len(ipStr)-1]), nil
}

/*
 * Port forwarding
 */
func (network *networkHandler) Forwarding(target string, port int) error {
	// Declare listener to the origin port
	listener, err := net.Listen("tcp", portStringify(port))
	if err != nil {
		return err
	}

	for {
		// Start connection to listener
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		fmt.Printf("Accepted connection %v\n", conn)
		go network.forward(conn, target)
	}
}

/*
 * Forward connection to then given port number
 */
func (network *networkHandler) forward(conn net.Conn, target string) {
	// Declare client to the forwarding port
	client, err := net.Dial("tcp", target)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to localhost %v\n", conn)

	// Copy IO
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

/*
 * Make port number to string format ":xxx"
 */
func portStringify(port int) string {
	return fmt.Sprintf(":%v", port)
}
