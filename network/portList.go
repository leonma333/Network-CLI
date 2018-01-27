package main

import (
  "fmt"
  "strings"
  "strconv"
)

// Create a new type for a list of port number
type portList []int

/*
 * Convert from portList to string implementation
 */
func(pl *portList) String() string {
  return fmt.Sprintf("%v", *pl)
}

/*
 * Convert from string to portList implementation
 */
func (pl *portList) Set(value string) error {
  convertStringArrayToPortlist := func(strArr []string) portList {
    var list portList
    for _, str := range strArr {
      port, _ := strconv.Atoi(str)
      list = append(list, port)
    }
    return list
  }
  *pl = convertStringArrayToPortlist(strings.Split(value, ","))
  return nil
}
