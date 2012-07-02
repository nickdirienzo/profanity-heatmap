package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Event struct {
	Payload Payload `json:"payload"`
	Actor   string  `json:"actor"`
	Type    string  `json:"type"`
}

type Payload struct {
	Shas []interface{} `json:"shas"`
}

type User struct {
	Location string `json:"location"`
}

const GITHUB string = "https://api.github.com"
const GEOCODER string = "http://maps.googleapis.com/maps/api/geocode/json?address="

func (e *Event) GetActorLocation() (string, error) {
	param := fmt.Sprint("/users/%q", e.Actor)
	resp, err := http.Get(GITHUB + param)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	var user User
	for {
		if err := decoder.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}
	return user.Location, nil
}

func GetTimeline(year int, month int, day int, hour int) (*Reader, error) {
	url := fmt.Sprint("http://data.githubarchive.org/%d-%d-%d-%d.json.gz", year, month, day, hour)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func ConvertToLatLng(loc string) (lat float32, lng float32) {
	param := strings.Join(strings.Split(loc, " "), "+")
	url := (GEOCODER + param + "&sensor=false")
	resp, err := http.Get(url)
	//TODO: Finish this function
}

func main() {
	file, err := os.Open("2012-03-11-15.json.gz")
	if err != nil {
		return err
	}
	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(reader)
	var event Event
	for i := 0; i < 5; i++ {
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Couldn't decode JSON:", err)
		}

		if event.Type == "PushEvent" {
			for _, commit := range event.Payload.Shas {
				switch commitData := commit.(type) {
				case []interface{}:
					fmt.Println("Commit Message:", commitData[2])
					loc, err := event.GetActorLocation()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(event.Actor, loc, ConvertToLatLng(loc))
				default:
					log.Fatal("Could not convert Shas to something usable")
				}
			}
		}

	}
}
