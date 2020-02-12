package schedule_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rltvty/go-home/timeline/schedule"
)

var _ = Describe("Schedule", func() {
	Describe("ParseJson", func() {
		var inJSON string
		Context("with valid json", func() {
			BeforeEach(func() {
				inJSON = `[{
					"id": 1,
					"channel": "twitch.tv/ninja",
					"type": "GUESSTIMATE",
					"startsAt": "Thu Jan 17 10:12:00 PST 2019",
					"endsAt": "Thu Jan 17 16:21:00 PST 2019",
					"priority": 10
				  }, {
					"id": 2,
					"channel": "twitch.tv/lirik",
					"type": "SCHEDULED",
					"startsAt": "Thu Jan 17 08:02:00 PST 2019",
					"endsAt": "Thu Jan 17 13:12:00 PST 2019",
					"priority": 7
				  }]`
			})
			It("should return an slice of 2 Stream objects, and not error", func() {
				streams, err := ParseJSON([]byte(inJSON))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(*streams)).To(Equal(2))

				Expect((*streams)[0].ID).To(Equal(1))
				Expect((*streams)[0].Channel).To(Equal("twitch.tv/ninja"))
				Expect((*streams)[0].StreamType).To(Equal(GUESSTIMATE))
				Expect((*streams)[0].StartsAt.Format(time.Stamp)).To(Equal("Jan 17 10:12:00"))
				Expect((*streams)[0].Priority).To(Equal(10))

				Expect((*streams)[1].ID).To(Equal(2))
				Expect((*streams)[1].Channel).To(Equal("twitch.tv/lirik"))
				Expect((*streams)[1].StreamType).To(Equal(SCHEDULED))
				Expect((*streams)[1].EndsAt.Format(time.Stamp)).To(Equal("Jan 17 13:12:00"))
				Expect((*streams)[1].Priority).To(Equal(7))
			})
		})
	})

	Describe("GetStartAndEnd", func() {
		Context("when there are two streams", func() {
			streamJSON := `[{
				"id": 1,
				"channel": "twitch.tv/ninja",
				"type": "GUESSTIMATE",
				"startsAt": "Thu Jan 17 10:12:00 PST 2019",
				"endsAt": "Thu Jan 17 16:21:00 PST 2019",
				"priority": 10
			  },{
				"id": 7,
				"channel": "twitch.tv/shroud",
				"type": "GUESSTIMATE",
				"startsAt": "Thu Jan 17 18:52:00 PST 2019",
				"endsAt": "Thu Jan 17 22:41:00 PST 2019",
				"priority": 0
			  }]`
			It("should return times that span the whole schedule", func() {
				streams, err := ParseJSON([]byte(streamJSON))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(*streams)).To(Equal(2))
				startsAt, endsAt := GetStartAndEnd(*streams)
				Expect(startsAt.Format(time.Stamp)).To(Equal("Jan 17 10:12:00"))
				Expect(endsAt.Format(time.Stamp)).To(Equal("Jan 17 22:41:00"))
			})
		})
	})

	Describe("GetPotentialStreams", func() {
		Context("when 3 streams exist, but 2 are in the window", func() {
			It("should return 2 streams, sorted by priority", func() {
				streams := []Stream{
					{
						ID:         1,
						Channel:    "test/1",
						StreamType: SCHEDULED,
						Priority:   5,
						StartsAt:   time.Now().Add(2 * time.Hour),
						EndsAt:     time.Now().Add(time.Hour),
					},
					{
						ID:         2,
						Channel:    "test/2",
						StreamType: SCHEDULED,
						Priority:   2,
						StartsAt:   time.Now(),
						EndsAt:     time.Now().Add(time.Hour),
					},
					{
						ID:         3,
						Channel:    "test/3",
						StreamType: SCHEDULED,
						Priority:   8,
						StartsAt:   time.Now(),
						EndsAt:     time.Now().Add(time.Hour),
					},
				}
				potentialStreams := GetPotentialStreams(time.Now().Add(30*time.Minute), streams)
				Expect(len(streams)).To(Equal(3))
				Expect(len(potentialStreams)).To(Equal(2))
				Expect(potentialStreams[0].ID).To(Equal(3))
				Expect(potentialStreams[1].ID).To(Equal(2))
			})
		})
	})

	Describe("GetSchedule", func() {
		Context("when there are no streams", func() {
			It("should return an empty schedule", func() {
				sched := GetSchedule([]Stream{})
				Expect(len(sched)).To(Equal(0))
			})
		})

		Context("when there is one stream", func() {
			It("should return a schedule with the only stream", func() {
				stream := Stream{
					ID:         1,
					Channel:    "test/1",
					StreamType: SCHEDULED,
					Priority:   5,
					StartsAt:   time.Now(),
					EndsAt:     time.Now().Add(time.Hour),
				}

				sched := GetSchedule([]Stream{stream})
				Expect(len(sched)).To(Equal(1))
				Expect(sched[0].Channel).To(Equal("test/1"))
			})
		})
	})

})
