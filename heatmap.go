package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	//"net/http"
	"os"
)

type Event struct {
	Payload Payload `json:"payload"`
	Type  string  `json:"type"`
}

type Payload struct {
    Shas []interface{} `json:"shas"`
}

func main() {
	/*resp, err := http.Get("http://data.githubarchive.org/2012-03-11-12.json.gz")
	if err != nil {
		fmt.Println("Couldn't retrieve data from GitHub Archive")
		return
	}
	defer resp.Body.Close()
	reader, err := gzip.NewReader(resp.Body)*/
	file, err := os.Open("2012-03-11-15.json.gz")
	if err != nil {
		fmt.Println("Couldn't open file")
		return
	}
	reader, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println("Couldn't uncompress")
		return
	}
	decoder := json.NewDecoder(reader)
    var event Event
	for i := 0; i < 10; i++ {
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
                    default:
                        log.Fatal("Could not convert Shas to something usable")
                }
            }
        }

	}
}
