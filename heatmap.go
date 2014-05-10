package main

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"log"
	"time"
)

const (
	dbName     string = "heatmap"
	eventsName string = "events"
)

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
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(DB(session))
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})
	go getDailyActivity(session)
	m.Run()
}

func getDailyActivity(s *mgo.Session) {
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
		for i := range events {
			log.Printf("Inserting event: %v", i)
			err = s.DB(dbName).C(eventsName).Insert(events[i])
			if err != nil {
				log.Printf("Insert error: %v", err)
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
