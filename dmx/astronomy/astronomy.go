package astronomy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const astronomyURL = "https://api.sunrise-sunset.org"
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
		URL: astronomyURL,
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

func (e Events) String() string  {
	dawn := e.Dawn.Local().Format("15:04")
	rise := e.SunRise.Local().Format("15:04")
	peak := e.SunPeak.Local().Format("15:04")
	set := e.SunSet.Local().Format("15:04")
	dusk := e.Dusk.Local().Format("15:04")

	return fmt.Sprintf("Dawn: %s   Rise: %s   Peak: %s   Set: %s   Dusk: %s", dawn, rise, peak, set, dusk)
}

type apiEvents struct {
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
}

type apiResponse struct {
	Results json.RawMessage `json:"results"`
	Status  string          `json:"status"`
}

var statusMap = map[string]string{
	"OK":              "success",
	"INVALID_REQUEST": "either the lat or lng parameters are missing or invalid",
	"INVALID_DATE":    "the date parameter is missing or invalid",
	"UNKNOWN_ERROR":   "the request could not be processed due to a server error. the request may succeed if you try again",
}

//GetEvents returns astronomical event times
func (api *API) GetEvents() (*Events, error) {
	path := fmt.Sprintf("/json?lat=%f&lng=%f&formatted=0", myLatitude, myLongitude)
	resp, err := api.Client.Get(fmt.Sprintf("%s%s", api.URL, path))
	if err != nil {
		return nil, fmt.Errorf("Error making api request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received unexpected response code: %s", resp.Status)
	}

	var response *apiResponse
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal the response json: %s", err)
	}

	if response.Status != "OK" {
		fmt.Println(response)
		return nil, fmt.Errorf("%s: %s", response.Status, statusMap[response.Status])
	}

	var results *apiEvents
	err = json.Unmarshal(response.Results, &results)
	if err != nil {
		return nil, err
	}
	if results == nil {
		return nil, errors.New("Results json was empty")
	}

	return &Events{
		Dawn:    results.CivilTwilightBegin,
		SunRise: results.Sunrise,
		SunPeak: results.SolarNoon,
		SunSet:  results.Sunset,
		Dusk:    results.CivilTwilightEnd,
	}, nil
}
