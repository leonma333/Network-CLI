package network_test

import (
	"testing"

	"."

	"github.com/stretchr/testify/assert"
)

type testCasePortList struct {
	list            network.PortList
	listString      string
	listStringArray string
}

var testCasesPortList = []testCasePortList{
	{
		list:            network.PortList{},
		listString:      "[]",
		listStringArray: "",
	},
	{
		list:            network.PortList{0},
		listString:      "[0]",
		listStringArray: "0",
	},
	{
		list:            network.PortList{3000},
		listString:      "[3000]",
		listStringArray: "3000",
	},
	{
		list:            network.PortList{3000, 4000, 5000},
		listString:      "[3000 4000 5000]",
		listStringArray: "3000,4000,5000",
	},
	{
		list:            network.PortList{},
		listString:      "[]",
		listStringArray: "love",
	},
	{
		list:            network.PortList{8080},
		listString:      "[8080]",
		listStringArray: "love,8080",
	},
}

func TestPortListString(t *testing.T) {
	for _, test := range testCasesPortList {
		assert.Equal(t, test.listString, test.list.String())
	}
}

func TestPortListSet(t *testing.T) {
	for _, test := range testCasesPortList {
		var portList network.PortList
		err := portList.Set(test.listStringArray)
		assert.NoError(t, err)
		assert.Equal(t, test.list, portList)
	}
}
