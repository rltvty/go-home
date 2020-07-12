package connection_test

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"

	. "github.com/rltvty/go-home/presonus/connection"
)

//{"speaker": {"Source":"locator","Mode":"broadcast","Port":53781,"Model":"SL328AI SPK","MacAddress":"00:0A:92:D6:66:EE","Kind":"speaker"}}"10.10.10.230"

var _ = Describe("Connection", func() {
	Describe("Connect", func() {
		StartManager()
	})
})
