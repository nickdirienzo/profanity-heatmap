package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

// To retrieve activity for April 11, 2012 at 3PM PST, %s would be substituted with 2012-04-11-15
const github_base_url string = "http://data.githubarchive.org/%s.json.gz"

// This regexp can be better
var profanityPattern *regexp.Regexp = regexp.MustCompile("([^\\w]ass[^\\w]|[^\\w]asshole[^\\w]|[^\\w]dumbass[^\\w]|[^\\w]hell[^\\w]|fuck|shit|damn|bitch|bastard)+")

func formatDateForQuery(date time.Time) string {
	var m, d string

	month := date.Month()
	if month < 10 {
		m = fmt.Sprintf("0%d", month)
	} else {
		m = fmt.Sprintf("%d", month)
	}

	day := date.Day()
	if day < 10 {
		d = fmt.Sprintf("0%d", day)
	} else {
		d = fmt.Sprintf("%d", day)
	}

	return fmt.Sprintf("%d-%s-%s-%d", date.Year(), m, d, date.Hour())
}

func processJSON(reader io.Reader) []Event {
	var events []Event
	var event Event
	decoder := json.NewDecoder(reader)
	for {
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Could not decode event: %v", err)
		} else {
			if event.Type == "PushEvent" {
				// This is really annoying
				for _, commit := range event.Payload.Commits {
					switch commitData := commit.(type) {
					case []interface{}:
						message := commitData[2]
						switch messageData := message.(type) {
						case string:
							if profanityPattern.MatchString(messageData) {
								log.Println("Found a match: " + messageData)
								events = append(events, event)
							}
						}
					}
				}
			}
		}
	}
	log.Println("There are", len(events), "events.")
	return events
}

func getActivity(url string) ([]Event, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reader, err := decompress(resp.Body)
	if err != nil {
		log.Printf("Could not decompress response: %v", err)
		return nil, err
	}
	events := processJSON(reader)
	return events, err
}

func decompress(compressed io.Reader) (io.Reader, error) {
	return gzip.NewReader(compressed)
}

func GetActivityForDate(date time.Time) ([]Event, error) {
	query_date := formatDateForQuery(date)
	fmt.Println(query_date)
	query_url := fmt.Sprintf(github_base_url, query_date)
	fmt.Println(query_url)
	events, err := getActivity(query_url)
	return events, err
}
