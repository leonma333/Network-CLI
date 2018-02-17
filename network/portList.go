package network

import (
	"fmt"
	"strconv"
	"strings"
)

// Create a new type for a list of port number
type PortList []int

/*
 * Convert from PortList to string implementation
 */
func (pl *PortList) String() string {
	return fmt.Sprintf("%v", *pl)
}

/*
 * Convert from string to PortList implementation
 */
func (pl *PortList) Set(value string) error {
	convertStringArrayToPortlist := func(strArr []string) PortList {
		var list PortList
		for _, str := range strArr {
			port, _ := strconv.Atoi(str)
			list = append(list, port)
		}
		return list
	}
	*pl = convertStringArrayToPortlist(strings.Split(value, ","))
	return nil
}
