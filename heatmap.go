package main

import (
	"flag"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"math/rand"
	"time"
)

var apiKey = flag.String("key", "", "Google Geocoding API Key")

const (
	dbName        string = "heatmap"
	eventsName    string = "events"
	eventsPerHour int    = 2500 / 24
)

type FrontendEvent struct {
	Lat      float64
	Lng      float64
	Location string
	Id       bson.ObjectId
	Count    int
}

// Adapted from: http://blog.gopheracademy.com/day-11-martini
func DB(session *mgo.Session) martini.Handler {
	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(dbName))
		defer s.Close()
		c.Next()
	}
}

func main() {
	flag.Parse()
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(DB(session))
	m.Use(martini.Static("static"))
	m.Get("/", func(r render.Render, db *mgo.Database) {
		var events []Event
		eventsMap := make(map[string][]Event)
		db.C(eventsName).Find(nil).All(&events)
		for i := range events {
			s := events[i].serialize()
			eventsMap[s] = append(eventsMap[s], events[i])
		}
		var frontendEvents []FrontendEvent
		for key := range eventsMap {
			mappedEvents := eventsMap[key]
			for i := range mappedEvents {
				frontendEvents = append(frontendEvents,
					FrontendEvent{mappedEvents[i].Lat, mappedEvents[i].Lng,
						mappedEvents[i].ActorAttributes.Location, mappedEvents[i].Id,
						len(mappedEvents)})
			}
		}
		r.HTML(200, "index", frontendEvents)
	})
	//go getDailyActivity(session)
	m.Run()
}

func getDailyActivity(s *mgo.Session) {
	g := Geocoder{*apiKey}
	for {
		session := s.Clone()
		now := time.Now()
		// I'm not sure when githubarchive.org updates,
		// and I don't feel like writing for BigQuery right now,
		// so I would hope they have data from 6 hours ago
		d, err := time.ParseDuration("-6h")
		if err != nil {
			log.Printf("Could not parse hour long duration: %v", err)
		}
		events, err := GetActivityForDate(now.Add(d))
		if err != nil {
			log.Printf("Could not get activity for date: %v", err)
		}
		p := rand.Perm(len(events))
		var event Event
		var eventsProcessed int = 0
		for i := 0; i < len(events); i++ {
			event = events[p[i]]
			if len(event.ActorAttributes.Location) > 0 {
				lat, lng, err := g.GetLatLng(event.ActorAttributes.Location)
				if err != nil && err.Error() != NO_ADDRESS {
					log.Printf("GetLatLng Error: %v", err)
				}
				event.Lat = lat
				event.Lng = lng
				err = s.DB(dbName).C(eventsName).Insert(event)
				if err != nil {
					log.Printf("Insert error: %v", err)
					eventsProcessed--
				}
				eventsProcessed++
				if eventsProcessed == eventsPerHour {
					break
				}
			}
		}
		d, err = time.ParseDuration("1h")
		if err != nil {
			log.Printf("Could not parse hour long duration: %v", err)
		}
		then := now.Add(d)
		log.Printf("Sleeping for %v", then.Sub(now))
		session.Close()
		time.Sleep(then.Sub(now))
	}
}
