package main

import (
	"log"
	"time"
)

func main() {
	pst, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatal("Could not create PST timezone: ", err.Error())
	}
	date := time.Date(2013, time.September, 4, 0, 0, 0, 0, pst)
	_ = GetActivityForDate(date)
}
