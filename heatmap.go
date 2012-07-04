package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Event struct {
	Payload   Payload `json:"payload"`
	Committer string  `json:"actor"`
	Type      string  `json:"type"`
}

type Payload struct {
	Shas []interface{} `json:"shas"`
}

type User struct {
	Location string `json:"location"`
}

type GeocodeResult struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float32
				Lng float32
			}
		}
	}
}

var (
	Commits          map[string]int                                                   // map of "lat,lng" hashed : count
	ProfanityPattern = regexp.MustCompile("\\w*(shit|piss|fuck)+\\w*|cunt|wtf|dafuq") // Used some of the words used here: http://goo.gl/Mvisn
	TotalCommits     int
)

const (
	GITHUB   string = "https://api.github.com"
	GEOCODER string = "http://maps.googleapis.com/maps/api/geocode/json?address="
)

func (e *Event) GetCommitterLocation() (string, error) {
	param := fmt.Sprint("/users/", e.Committer)
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

func (e *Event) GetLatLng() (float32, float32, error) {
	loc, err := e.GetCommitterLocation()
	if err != nil {
		return 0, 0, err
	}
	param := strings.Join(strings.Split(loc, " "), "+")
	url := (GEOCODER + param + "&sensor=false")
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	var gR GeocodeResult
	for {
		if err := decoder.Decode(&gR); err == io.EOF {
			break
		} else if err != nil {
			return 0, 0, err
		}
	}
	if len(gR.Results) < 1 {
		return 0, 0, errors.New("No Geocoder Results")
	}
	return gR.Results[0].Geometry.Location.Lat, gR.Results[0].Geometry.Location.Lng, nil
}

/*func GetTimelineReader(year int, month int, day int, hour int) (*Reader, error) {
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
}*/

func main() {
	Commits = make(map[string]int)
	file, err := os.Open("2012-03-11-15.json.gz")
	if err != nil {
		log.Fatal(err)
	}
	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(reader)
	var event Event
	for {
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			log.Print("Couldn't decode JSON:", err)
		}

		if event.Type == "PushEvent" {
			for _, commit := range event.Payload.Shas {
				switch commitData := commit.(type) {
				case []interface{}:
					message := commitData[2]
					switch m := message.(type) {
					case string:
						if ProfanityPattern.MatchString(m) {
							lat, lng, _ := event.GetLatLng()
							if lat != 0 && lng != 0 {
								key := fmt.Sprintf("%f,%f", lat, lng)
								count, _ := Commits[key]
								Commits[key] = (count + 1)
								fmt.Println(m, "by", event.Committer, "(Lat:", lat, "Lng", lng, ")")
							}
						}
					default:
						log.Fatal("Commit message was not a string")
					}
				default:
					log.Fatal("Could not convert Shas to something usable")
				}
			}
		}
		TotalCommits += 1
	}
	fmt.Println("Commits:", Commits, "Total Commits:", TotalCommits)
}
