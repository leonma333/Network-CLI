package network

import (
	"reflect"
	"testing"
)

var (
	portListTest            = PortList([]int{1000, 2000, 3000})
	portListTestString      = "[1000 2000 3000]"
	portListTestStringArray = "1000,2000,3000"
)

func TestPortListString(t *testing.T) {
	str := portListTest.String()
	if str != portListTestString {
		t.Errorf("Expect %s, but get %s", portListTestString, str)
	}
}

func TestPortListSet(t *testing.T) {
	list := portListTest
	list.Set(portListTestStringArray)
	if !reflect.DeepEqual(list, portListTest) {
		t.Errorf("Expect %v, but get %v", portListTest, list)
	}
}
