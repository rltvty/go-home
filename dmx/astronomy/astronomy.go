package astronomy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const astronomyAPIURL = "https://api.sunrise-sunset.org/json?lat=%f&lng=%f&formatted=0"
const myLatitude = 30.262890
const myLongitude = -97.720119

//API for accessing astronomy events
type API struct {
	Client *http.Client
	URL    string
}

//New creates the API client with optional options
func New(options ...func(*API)) *API {
	api := API{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		URL: fmt.Sprintf(astronomyAPIURL, myLatitude, myLongitude),
	}
	api.SetOptions(options...)

	return &api
}

// SetOptions takes one or more option function and applies them in order to API.
func (api *API) SetOptions(options ...func(*API)) {
	for _, opt := range options {
		opt(api)
	}
}

//Events contains time info about astronomical events
type Events struct {
	Dawn    time.Time
	SunRise time.Time
	SunPeak time.Time
	SunSet  time.Time
	Dusk    time.Time
}

type sunriseSunsetAPIResponse struct {
	Results *struct {
		Sunrise                   time.Time `json:"sunrise"`
		Sunset                    time.Time `json:"sunset"`
		SolarNoon                 time.Time `json:"solar_noon"`
		DayLength                 int       `json:"day_length"`
		CivilTwilightBegin        time.Time `json:"civil_twilight_begin"`
		CivilTwilightEnd          time.Time `json:"civil_twilight_end"`
		NauticalTwilightBegin     time.Time `json:"nautical_twilight_begin"`
		NauticalTwilightEnd       time.Time `json:"nautical_twilight_end"`
		AstronomicalTwilightBegin time.Time `json:"astronomical_twilight_begin"`
		AstronomicalTwilightEnd   time.Time `json:"astronomical_twilight_end"`
	} `json:"results,omitempty"`
	Status string `json:"status"`
}

var statusMap = map[string]string{
	"OK":              "success",
	"INVALID_REQUEST": "either the lat or lng parameters are missing or invalid",
	"INVALID_DATE":    "the date parameter is missing or invalid",
	"UNKNOWN_ERROR":   "the request could not be processed due to a server error. the request may succeed if you try again",
}

//GetEvents returns astronomical event times
func (api *API) GetEvents() (*Events, error) {
	resp, err := api.Client.Get(api.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response sunriseSunsetAPIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	if response.Status != "OK" {
		return nil, fmt.Errorf("%s: %s", response.Status, statusMap[response.Status])
	}

	if response.Results == nil {
		return nil, errors.New("MISSING_RESULTS: request was okay, but results are missing. try again")
	}

	results := response.Results
	return &Events{
		Dawn:    results.CivilTwilightBegin,
		SunRise: results.Sunrise,
		SunPeak: results.SolarNoon,
		SunSet:  results.Sunset,
		Dusk:    results.CivilTwilightEnd,
	}, nil
}
