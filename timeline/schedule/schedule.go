package schedule

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type streamType int

const (
	GUESSTIMATE streamType = 0 + iota
	SCHEDULED
)

//Stream defines the properties of a stream
type Stream struct {
	ID         int
	Channel    string
	StreamType streamType
	StartsAt   time.Time
	EndsAt     time.Time
	Priority   int
}

type byPriorityAndType []Stream

//Implement sort interface on byPriorityAndType
func (s byPriorityAndType) Len() int {
	return len(s)

}
func (s byPriorityAndType) Less(i, j int) bool {
	if s[i].Priority > s[j].Priority {
		return true
	}
	return s[i].StreamType < s[j].StreamType
}

func (s byPriorityAndType) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//StreamJSON intermediate type for parsing the json data
type StreamJSON struct {
	ID       int    `json:"id"`
	Channel  string `json:"channel"`
	TextType string `json:"type"`
	StartsAt string `json:"startsAt"`
	EndsAt   string `json:"endsAt"`
	Priority int    `json:"priority"`
}

//Item an item on the outputted schedule
type Item struct {
	Channel  string    `json:"channel"`
	StreamID int       `json:"streamId"`
	Stream   Stream    `json:"-"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

//ParseJSON parses the input data into a slice of strongly typed Stream objects
func ParseJSON(jsonData []byte) (*[]Stream, error) {
	var jsonStreams []StreamJSON
	err := json.Unmarshal(jsonData, &jsonStreams)
	if err != nil {
		return nil, err
	}

	streams := []Stream{}
	for _, jsonStream := range jsonStreams {

		var sType streamType
		switch jsonStream.TextType {
		case "GUESSTIMATE":
			sType = GUESSTIMATE
		case "SCHEDULED":
			sType = SCHEDULED
		default:
			return nil, fmt.Errorf("Unknown StreamType: %s", jsonStream.TextType)
		}

		startsAt, err := time.Parse(time.UnixDate, jsonStream.StartsAt)
		if err != nil {
			return nil, err
		}
		endsAt, err := time.Parse(time.UnixDate, jsonStream.EndsAt)
		if err != nil {
			return nil, err
		}

		streams = append(streams, Stream{
			ID:         jsonStream.ID,
			Channel:    jsonStream.Channel,
			StreamType: sType,
			StartsAt:   startsAt,
			EndsAt:     endsAt,
			Priority:   jsonStream.Priority,
		})
	}

	return &streams, nil
}

//GetStartAndEnd gets the start and end times for the schedule
func GetStartAndEnd(streams []Stream) (*time.Time, *time.Time) {
	if len(streams) == 0 {
		return nil, nil
	}
	startsAt := streams[0].StartsAt
	endsAt := streams[0].EndsAt
	for _, stream := range streams[1:] {
		if stream.StartsAt.Before(startsAt) {
			startsAt = stream.StartsAt
		}
		if stream.EndsAt.After(endsAt) {
			endsAt = stream.EndsAt
		}
	}
	return &startsAt, &endsAt
}

//GetPotentialStreams gets all the streams that are available at the currentTime, ordered by priority
func GetPotentialStreams(currentTime time.Time, streams []Stream) []Stream {
	potentialStreams := []Stream{}
	for _, stream := range streams {
		if stream.StartsAt.Before(currentTime) && stream.EndsAt.After(currentTime) {
			potentialStreams = append(potentialStreams, stream)
		}
	}
	sort.Sort(byPriorityAndType(potentialStreams))
	return potentialStreams
}

func shouldSwitchStreams(currentStream *Stream, streamChoices []Stream, currentTime time.Time, watchedIds map[int]struct{}) *Stream {
	if (currentStream != nil && currentStream.Priority == 10) || len(streamChoices) == 0 {
		return nil
	}

	alreadyWatched := func(choice Stream) bool {
		_, ok := watchedIds[choice.ID]
		return ok
	}

	//filter alreadyWatched Streams, unless priority 10
	filteredChoices := []Stream{}
	for _, choice := range streamChoices {
		if choice.Priority == 10 || !alreadyWatched(choice) {
			filteredChoices = append(filteredChoices, choice)
		}
	}

	//if no un-watched choices are avialable, reset to initial list
	if len(filteredChoices) == 0 {
		filteredChoices = streamChoices
	}

	//filter guestamite type if it didn't start at least 10 mins ago, unless priority 10
	filteredChoices2 := []Stream{}
	for _, choice := range filteredChoices {
		if choice.Priority == 10 || choice.StreamType == SCHEDULED || currentTime.After(choice.StartsAt.Add(10*time.Minute)) {
			filteredChoices2 = append(filteredChoices2, choice)
		}
	}

	//if no choices are avialable, reset to previous list
	if len(filteredChoices2) == 0 {
		filteredChoices2 = filteredChoices
	}

	//Grab the highest priority choice from the filtered set
	bestChoice := filteredChoices2[0]

	if currentStream == nil || currentStream.Priority+4 < bestChoice.Priority {
		return &bestChoice
	}
	return nil
}

//GetSchedule returns the recommened watching schedule based on the available stream choices
func GetSchedule(streams []Stream) []Item {
	items := []Item{}

	startAt, endsAt := GetStartAndEnd(streams)
	if startAt == nil || endsAt == nil {
		return items
	}
	currentTime := *startAt
	var currentItem *Item = nil
	alreadyWatchedStreamIDs := map[int]struct{}{}

	saveCurrentItem := func() {
		currentItem.EndsAt = currentTime.Add(-1 * time.Minute)
		items = append(items, *currentItem)
		currentItem = nil
	}

	for {
		if currentTime.After(*endsAt) {
			break
		}

		potentialStreams := GetPotentialStreams(currentTime, streams)
		var currentStream *Stream = nil
		if currentItem != nil {
			currentStream = &currentItem.Stream
		}
		nextStream := shouldSwitchStreams(currentStream, potentialStreams, currentTime, alreadyWatchedStreamIDs)
		if nextStream != nil {
			if currentItem != nil {
				saveCurrentItem()
			}
			currentItem = &Item{
				StreamID: nextStream.ID,
				Stream:   *nextStream,
				StartsAt: currentTime.Add(-1 * time.Minute),
				Channel:  nextStream.Channel,
			}
		} else if currentItem != nil && currentTime.After(currentItem.Stream.EndsAt) {
			saveCurrentItem()
		}

		currentTime = currentTime.Add(time.Minute)
	}

	saveCurrentItem()
	return items
}
