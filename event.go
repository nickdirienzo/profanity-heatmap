package main

import "labix.org/v2/mgo/bson"

type Payload struct {
	// ["hash", "email", "message", "name", distinct_bool]
	Commits []interface{} `json:"shas"`
}

type ActorAttributes struct {
	Login      string `json:"login" bson:"login"`
	Type       string `json:"type" bson:"type"`
	GravatarId string `json:"gravatar_id" bson:"gravatar_id"`
	Name       string `json:"name" bson:"name"`
	Location   string `json:"location" bson:"location"`
}

type Event struct {
	Id              bson.ObjectId   `bson:"_id,omitempty"`
	Payload         Payload         `json:"payload" bson:"-"`
	Actor           string          `json:"actor" bson:"actor"`
	ActorAttributes ActorAttributes `json:"actor_attributes" bson:"actor_attributes,inline"`
	Type            string          `json:"type" bson:"event_type"`
	Lat             float64         `bson:"lat"`
	Lng             float64         `bson:"lng"`
}
