package astronomy_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/rltvty/go-home/dmx/astronomy"
)

var _ = Describe("GetEvents", func() {
	var server *ghttp.Server
	var client *API

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = New(func(api *API) {
			api.URL = fmt.Sprintf("%s/json", server.URL())
		})
	})

	AfterEach(func() {
		//shut down the server between tests
		server.Close()
	})

	Describe("fetching json", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.VerifyRequest("GET", "/json"),
			)
		})

		It("should make a request to the json endpoint", func() {
			client.GetEvents()
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})
})
