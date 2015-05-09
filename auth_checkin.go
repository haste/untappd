package untappd

import (
	"net/http"
	"net/url"
	"strconv"
)

type CheckinRequest struct {
	// Mandatory parameters
	BeerID    int
	GMTOffset int
	TimeZone  string

	// Optional parameters

	// Checkin location
	FoursquareID string
	Latitude     float64
	Longitude    float64

	// User comment and rating
	Comment string
	Rating  float64

	// Send to social media?
	Facebook bool
	Twitter  bool
	// FoursquareID is required if this is true
	Foursquare bool
}

func (a *AuthService) Checkin(r CheckinRequest) (*Checkin, *http.Response, error) {
	// Add required parameters
	q := url.Values{
		"bid":        []string{strconv.Itoa(r.BeerID)},
		"gmt_offset": []string{strconv.Itoa(r.GMTOffset)},
		"timezone":   []string{r.TimeZone},
	}

	// Add optional parameters, if not empty
	if r.FoursquareID != "" {
		q.Set("foursquare_id", r.FoursquareID)
	}
	if r.Latitude != 0 {
		q.Set("geolat", strconv.FormatFloat(r.Latitude, 'f', -1, 64))
	}
	if r.Longitude != 0 {
		q.Set("geolng", strconv.FormatFloat(r.Longitude, 'f', -1, 64))
	}

	if r.Comment != "" {
		q.Set("shout", r.Comment)
	}
	if r.Rating != 0 {
		q.Set("rating", strconv.FormatFloat(r.Rating, 'f', -1, 64))
	}

	if r.Facebook {
		q.Set("facebook", "on")
	}
	if r.Twitter {
		q.Set("twitter", "on")
	}
	if r.Foursquare {
		q.Set("foursquare", "on")
	}

	// Temporary struct to unmarshal checkin JSON
	var v struct {
		Response rawCheckin `json:"response"`
	}

	// Perform request to check in a beer
	res, err := a.client.request("POST", "checkin/add", q, nil, &v)
	if err != nil {
		return nil, res, err
	}

	return v.Response.export(), res, nil
}
