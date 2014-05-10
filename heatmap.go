package main

import (
	"log"
	"time"
)

func main() {
	go getDailyActivity()
	select {}
}

func getDailyActivity() {
	for {
		now := time.Now()
		// githubarchive.org only seems to have data from 4 hours ago and before
		d, err := time.ParseDuration("-4h")
		if err != nil {
			log.Printf("Could not parse hour long duration: %v", err)
		}
		_, err = GetActivityForDate(now.Add(d))
		d, err = time.ParseDuration("1h")
		if err != nil {
			log.Printf("Could not parse hour long duration: %v", err)
		}
		then := now.Add(d)
		log.Printf("Sleeping for %v", then.Sub(now))
		time.Sleep(then.Sub(now))
	}
}
