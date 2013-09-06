package main

type Payload struct {
	// ["hash", "email", "message", "name", distinct_bool]
	Shas []interface{} `json:"shas"`
}

type ActorAttributes struct {
	Login      string `json:"login"`
	Type       string `json:"type"`
	GravatarId string `json:"gravatar_id"`
	Name       string `json:"name"`
	Locatoin   string `json:"location"`
}

type Event struct {
	Payload         Payload         `json:"payload"`
	Actor           string          `json:"actor"`
	ActorAttributes ActorAttributes `json:"actor_attributes"`
	Type            string          `json:"type"`
}
