package astronomy_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/rltvty/go-home/dmx/astronomy"
)

func parseTimeHelper(input string) time.Time {
	output, _ := time.Parse(time.RFC3339, input)
	return output
}

var _ = Describe("GetEvents", func() {
	var server *ghttp.Server
	var client *API
	var statusCode int
	var responseBody string

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = New(func(api *API) {
			api.URL = server.URL()
		})
	})

	AfterEach(func() {
		//shut down the server between tests
		server.Close()
	})

	Describe("fetching json", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/json", "lat=30.262890&lng=-97.720119&formatted=0"),
					ghttp.RespondWithPtr(&statusCode, &responseBody),
				),
			)
		})

		Context("when the request succeeds", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				responseBody = `{
					"results":
					{
					  "sunrise":"2015-05-21T05:05:35+00:00",
					  "sunset":"2015-05-21T19:22:59+00:00",
					  "solar_noon":"2015-05-21T12:14:17+00:00",
					  "day_length":51444,
					  "civil_twilight_begin":"2015-05-21T04:36:17+00:00",
					  "civil_twilight_end":"2015-05-21T19:52:17+00:00",
					  "nautical_twilight_begin":"2015-05-21T04:00:13+00:00",
					  "nautical_twilight_end":"2015-05-21T20:28:21+00:00",
					  "astronomical_twilight_begin":"2015-05-21T03:20:49+00:00",
					  "astronomical_twilight_end":"2015-05-21T21:07:45+00:00"
					},
					 "status":"OK"
				  }`
			})

			It("should return the event times without erroring", func() {
				events, err := client.GetEvents()
				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(HaveOccurred())
				Expect(events).To(Equal(&Events{
					Dawn:    parseTimeHelper("2015-05-21T04:36:17+00:00"),
					SunRise: parseTimeHelper("2015-05-21T05:05:35+00:00"),
					SunPeak: parseTimeHelper("2015-05-21T12:14:17+00:00"),
					SunSet:  parseTimeHelper("2015-05-21T19:22:59+00:00"),
					Dusk:    parseTimeHelper("2015-05-21T19:52:17+00:00"),
				}))
			})
		})

		Context("when the server doesn't exist", func() {
			BeforeEach(func() {
				client = New(func(api *API) {
					api.URL = "http://nope"
					api.Client = &http.Client{
						Timeout: 2 * time.Second,
					}
				})
			})

			It("should return an error", func() {
				events, err := client.GetEvents()
				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).To(MatchError("Error making api request: Get http://nope/json?lat=30.262890&lng=-97.720119&formatted=0: dial tcp: lookup nope: no such host"))
				Expect(events).To(BeNil())
			})
		})

		Context("when the endpoint is not found", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
				responseBody = ""
			})

			It("should return an error", func() {
				events, err := client.GetEvents()
				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(MatchError("Received unexpected response code: 404 Not Found"))
				Expect(events).To(BeNil())
			})
		})

		Context("when the request returns an error status", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				responseBody = `{"results":"","status":"INVALID_REQUEST"}`
			})

			It("should return an error", func() {
				events, err := client.GetEvents()
				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(MatchError("INVALID_REQUEST: either the lat or lng parameters are missing or invalid"))
				Expect(events).To(BeNil())
			})
		})

		Context("when the request returns invalid json", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				responseBody = `{"result":"","status":"OK"}`
			})

			It("should return an error", func() {
				events, err := client.GetEvents()
				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(MatchError("Unable to unmarshal the response json: json: unknown field \"result\""))
				Expect(events).To(BeNil())
			})
		})

		Context("when the request returns null results json", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				responseBody = `{"results":null,"status":"OK"}`
			})

			It("should return an error", func() {
				events, err := client.GetEvents()
				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(MatchError("Results json was empty"))
				Expect(events).To(BeNil())
			})
		})

		Context("when the request returns invalid results json", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				responseBody = `{"results":"","status":"OK"}`
			})

			It("should return an error", func() {
				events, err := client.GetEvents()
				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).To(MatchError("json: cannot unmarshal string into Go value of type astronomy.apiEvents"))
				Expect(events).To(BeNil())
			})
		})
	})
})
