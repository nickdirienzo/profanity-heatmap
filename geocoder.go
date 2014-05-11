package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const google_base_url string = "https://maps.googleapis.com/maps/api/geocode/json"
const NO_ADDRESS string = "No address provided"

type Geocoder struct {
	apiKey string
}

type GeocodingResults struct {
	Results []GeocodingResult `json:"results"`
}

type GeocodingResult struct {
	Geometry GeocodingGeometry `json:"geometry"`
}

type GeocodingGeometry struct {
	Location GeocodingLocation `json:"location"`
}

type GeocodingLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (g *Geocoder) GetLatLng(address string) (float64, float64, error) {
	if len(address) > 0 {
		req, err := url.Parse(google_base_url)
		if err != nil {
			log.Fatal("Could not generate request from base_url: %v", err)
		}
		values := req.Query()
		values.Add("address", address)
		values.Add("sensor", "false")
		values.Add("key", g.apiKey)
		req.RawQuery = values.Encode()
		log.Println("Making request for: " + req.String())
		resp, err := http.Get(req.String())
		if err != nil {
			log.Printf("Could not make request: %v", err)
			return -1, -1, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Could not read response: %v", err)
			return -1, -1, err
		}
		var results GeocodingResults
		err = json.Unmarshal(body, &results)
		if err != nil {
			log.Printf("Could not create GeocodingResults: %v", err)
			return -1, -1, err
		}
		return 0, 0, nil
	}
	return -1, -1, errors.New(NO_ADDRESS)
}
